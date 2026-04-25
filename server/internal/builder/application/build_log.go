package application

import (
	"bytes"
	"context"
	"io"
	"log"
	"sync"

	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
)

type BuildLogApplication struct {
	mu      sync.Mutex
	seq     uint64
	buffers map[int64]*bytes.Buffer
	subs    map[int64]map[uint64]chan *value.BuildLogChunk
}

var _ port.LogPublisher = (*BuildLogApplication)(nil)

func NewBuildLogApplication() *BuildLogApplication {
	return &BuildLogApplication{
		buffers: make(map[int64]*bytes.Buffer),
		subs:    make(map[int64]map[uint64]chan *value.BuildLogChunk),
	}
}

func (a *BuildLogApplication) StreamBuildLogs(ctx context.Context, buildId int64) (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	ch, cancel, snapshot := a.subscribeWithSnapshot(buildId)
	go func() {
		defer cancel()
		closePW := func(err error) {
			if err != nil {
				_ = pw.CloseWithError(err)
			} else {
				_ = pw.Close()
			}
		}
		if len(snapshot) > 0 {
			if _, err := pw.Write(snapshot); err != nil {
				closePW(err)
				return
			}
		}
		for {
			select {
			case <-ctx.Done():
				closePW(ctx.Err())
				return
			case c, ok := <-ch:
				if !ok {
					closePW(nil)
					return
				}
				if c.End {
					closePW(nil)
					return
				}
				if len(c.Data) == 0 {
					continue
				}
				if _, err := pw.Write(c.Data); err != nil {
					closePW(err)
					return
				}
			}
		}
	}()
	return pr, nil
}

func (a *BuildLogApplication) subscribeWithSnapshot(buildId int64) (ch <-chan *value.BuildLogChunk, cancel func(), snapshot []byte) {
	const buf = 256
	typedCh := make(chan *value.BuildLogChunk, buf)
	a.mu.Lock()
	id := a.seq
	a.seq++
	if a.subs[buildId] == nil {
		a.subs[buildId] = make(map[uint64]chan *value.BuildLogChunk)
	}
	a.subs[buildId][id] = typedCh
	if b, ok := a.buffers[buildId]; ok && b.Len() > 0 {
		snapshot = make([]byte, b.Len())
		copy(snapshot, b.Bytes())
	}
	a.mu.Unlock()

	cancel = func() {
		a.mu.Lock()
		if m, ok := a.subs[buildId]; ok {
			if c, ok := m[id]; ok {
				delete(m, id)
				if len(m) == 0 {
					delete(a.subs, buildId)
				}
				close(c)
			}
		}
		a.mu.Unlock()
	}
	return typedCh, cancel, snapshot
}

func (a *BuildLogApplication) PublishLogChunk(buildId int64, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	a.mu.Lock()
	if a.buffers[buildId] == nil {
		a.buffers[buildId] = new(bytes.Buffer)
	}
	_, _ = a.buffers[buildId].Write(data)
	m := a.subs[buildId]
	chs := make([]chan *value.BuildLogChunk, 0, len(m))
	for _, c := range m {
		chs = append(chs, c)
	}
	a.mu.Unlock()

	for _, c := range chs {
		p := make([]byte, len(data))
		copy(p, data)
		chunk := &value.BuildLogChunk{BuildId: buildId, Data: p}
		select {
		case c <- chunk:
		default:
			log.Printf("build log stream buffer full for build %d, dropping chunk", buildId)
		}
	}
	return nil
}

func (a *BuildLogApplication) PublishLogEnd(buildId int64) error {
	a.mu.Lock()
	delete(a.buffers, buildId)
	m, ok := a.subs[buildId]
	if !ok {
		a.mu.Unlock()
		return nil
	}
	chs := make([]chan *value.BuildLogChunk, 0, len(m))
	for _, c := range m {
		chs = append(chs, c)
	}
	a.mu.Unlock()

	end := &value.BuildLogChunk{BuildId: buildId, End: true}
	for _, c := range chs {
		select {
		case c <- end:
		default:
			log.Printf("build log end dropped for build %d (stream buffer full)", buildId)
		}
	}
	return nil
}
