package application

import (
	"bytes"
	"context"
	"io"
	"log"
	"sync"

	"starliner.app/internal/provisioner/domain/port"
	"starliner.app/internal/provisioner/domain/value"
)

type ClusterLogApplication struct {
	mu      sync.Mutex
	seq     uint64
	buffers map[int64]*bytes.Buffer
	subs    map[int64]map[uint64]chan *value.ClusterLogChunk
}

var _ port.LogPublisher = (*ClusterLogApplication)(nil)

func NewClusterLogApplication() *ClusterLogApplication {
	return &ClusterLogApplication{
		buffers: make(map[int64]*bytes.Buffer),
		subs:    make(map[int64]map[uint64]chan *value.ClusterLogChunk),
	}
}

func (a *ClusterLogApplication) StreamProvisioningLogs(ctx context.Context, clusterId int64) (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	ch, cancel, snapshot := a.subscribeWithSnapshot(clusterId)
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

func (a *ClusterLogApplication) subscribeWithSnapshot(clusterId int64) (ch <-chan *value.ClusterLogChunk, cancel func(), snapshot []byte) {
	const buf = 256
	typedCh := make(chan *value.ClusterLogChunk, buf)
	a.mu.Lock()
	id := a.seq
	a.seq++
	if a.subs[clusterId] == nil {
		a.subs[clusterId] = make(map[uint64]chan *value.ClusterLogChunk)
	}
	a.subs[clusterId][id] = typedCh
	if b, ok := a.buffers[clusterId]; ok && b.Len() > 0 {
		snapshot = make([]byte, b.Len())
		copy(snapshot, b.Bytes())
	}
	a.mu.Unlock()

	cancel = func() {
		a.mu.Lock()
		if m, ok := a.subs[clusterId]; ok {
			if c, ok := m[id]; ok {
				delete(m, id)
				if len(m) == 0 {
					delete(a.subs, clusterId)
				}
				close(c)
			}
		}
		a.mu.Unlock()
	}
	return typedCh, cancel, snapshot
}

func (a *ClusterLogApplication) PublishLogChunk(clusterId int64, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	a.mu.Lock()
	if a.buffers[clusterId] == nil {
		a.buffers[clusterId] = new(bytes.Buffer)
	}

	_, _ = a.buffers[clusterId].Write(data)
	m := a.subs[clusterId]
	chs := make([]chan *value.ClusterLogChunk, 0, len(m))
	for _, c := range m {
		chs = append(chs, c)
	}
	a.mu.Unlock()

	for _, c := range chs {
		p := make([]byte, len(data))
		copy(p, data)
		chunk := &value.ClusterLogChunk{
			ClusterId: clusterId,
			Data:      p,
		}
		select {
		case c <- chunk:
		default:
			log.Printf("log stream buffer full for cluster id %d, dropping chunk", clusterId)
		}
	}
	return nil
}

func (a *ClusterLogApplication) PublishLogEnd(clusterId int64) error {
	a.mu.Lock()
	delete(a.buffers, clusterId)
	m, ok := a.subs[clusterId]
	if !ok {
		a.mu.Unlock()
		return nil
	}
	chs := make([]chan *value.ClusterLogChunk, 0, len(m))
	for _, c := range m {
		chs = append(chs, c)
	}
	a.mu.Unlock()

	end := &value.ClusterLogChunk{ClusterId: clusterId, End: true}
	for _, c := range chs {
		select {
		case c <- end:
		default:
			log.Printf("log end dropped for cluster id %d (stream buffer full)", clusterId)
		}
	}
	return nil
}
