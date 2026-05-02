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
	buffers map[string]*bytes.Buffer
	subs    map[string]map[uint64]chan *value.ClusterLogChunk
}

var _ port.LogPublisher = (*ClusterLogApplication)(nil)

func NewClusterLogApplication() *ClusterLogApplication {
	return &ClusterLogApplication{
		buffers: make(map[string]*bytes.Buffer),
		subs:    make(map[string]map[uint64]chan *value.ClusterLogChunk),
	}
}

func (a *ClusterLogApplication) StreamProvisioningLogs(ctx context.Context, provisioningId string) (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	ch, cancel, snapshot := a.subscribeWithSnapshot(provisioningId)
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

func (a *ClusterLogApplication) subscribeWithSnapshot(provisioningId string) (ch <-chan *value.ClusterLogChunk, cancel func(), snapshot []byte) {
	const buf = 256
	typedCh := make(chan *value.ClusterLogChunk, buf)
	a.mu.Lock()
	id := a.seq
	a.seq++
	if a.subs[provisioningId] == nil {
		a.subs[provisioningId] = make(map[uint64]chan *value.ClusterLogChunk)
	}
	a.subs[provisioningId][id] = typedCh
	if b, ok := a.buffers[provisioningId]; ok && b.Len() > 0 {
		snapshot = make([]byte, b.Len())
		copy(snapshot, b.Bytes())
	}
	a.mu.Unlock()

	cancel = func() {
		a.mu.Lock()
		if m, ok := a.subs[provisioningId]; ok {
			if c, ok := m[id]; ok {
				delete(m, id)
				if len(m) == 0 {
					delete(a.subs, provisioningId)
				}
				close(c)
			}
		}
		a.mu.Unlock()
	}
	return typedCh, cancel, snapshot
}

func (a *ClusterLogApplication) PublishLogChunk(provisioningId string, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	a.mu.Lock()
	if a.buffers[provisioningId] == nil {
		a.buffers[provisioningId] = new(bytes.Buffer)
	}

	_, _ = a.buffers[provisioningId].Write(data)
	m := a.subs[provisioningId]
	chs := make([]chan *value.ClusterLogChunk, 0, len(m))
	for _, c := range m {
		chs = append(chs, c)
	}
	a.mu.Unlock()

	for _, c := range chs {
		p := make([]byte, len(data))
		copy(p, data)
		chunk := &value.ClusterLogChunk{
			ProvisioningId: provisioningId,
			Data:           p,
		}
		select {
		case c <- chunk:
		default:
			log.Printf("log stream buffer full for provisioning id %s, dropping chunk", provisioningId)
		}
	}
	return nil
}

func (a *ClusterLogApplication) PublishLogEnd(provisioningId string) error {
	a.mu.Lock()
	delete(a.buffers, provisioningId)
	m, ok := a.subs[provisioningId]
	if !ok {
		a.mu.Unlock()
		return nil
	}
	chs := make([]chan *value.ClusterLogChunk, 0, len(m))
	for _, c := range m {
		chs = append(chs, c)
	}
	a.mu.Unlock()

	end := &value.ClusterLogChunk{ProvisioningId: provisioningId, End: true}
	for _, c := range chs {
		select {
		case c <- end:
		default:
			log.Printf("log end dropped for provisioning id %s (stream buffer full)", provisioningId)
		}
	}
	return nil
}
