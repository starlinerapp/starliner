package sse

import (
	"encoding/json"
	"log"
	"sync"

	"starliner.app/internal/core/domain/value"
)

type environmentKey struct {
	correlationId string
	environmentId int64
}

type EnvironmentNotificationHub struct {
	mu      sync.RWMutex
	clients map[environmentKey]map[chan *value.EnvironmentNotification]struct{}
}

func NewEnvironmentNotificationHub() *EnvironmentNotificationHub {
	return &EnvironmentNotificationHub{
		clients: make(map[environmentKey]map[chan *value.EnvironmentNotification]struct{}),
	}
}

func (h *EnvironmentNotificationHub) Subscribe(correlationId string, environmentId int64) chan *value.EnvironmentNotification {
	key := environmentKey{correlationId: correlationId, environmentId: environmentId}
	ch := make(chan *value.EnvironmentNotification, 16)
	h.mu.Lock()
	if h.clients[key] == nil {
		h.clients[key] = make(map[chan *value.EnvironmentNotification]struct{})
	}
	h.clients[key][ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *EnvironmentNotificationHub) Unsubscribe(correlationId string, environmentId int64, ch chan *value.EnvironmentNotification) {
	key := environmentKey{correlationId: correlationId, environmentId: environmentId}
	h.mu.Lock()
	if subs, ok := h.clients[key]; ok {
		delete(subs, ch)
		if len(subs) == 0 {
			delete(h.clients, key)
		}
	}
	h.mu.Unlock()
	close(ch)
}

func (h *EnvironmentNotificationHub) Broadcast(correlationId string, environmentId int64, notification *value.EnvironmentNotification) {
	key := environmentKey{correlationId: correlationId, environmentId: environmentId}
	h.mu.RLock()
	subs := h.clients[key]
	h.mu.RUnlock()

	for ch := range subs {
		select {
		case ch <- notification:
		default:
			log.Printf("dropping notification for slow client on environment %d correlationId %s", environmentId, correlationId)
		}
	}
}

func (h *EnvironmentNotificationHub) WriteNotification(w *Writer, notification *value.EnvironmentNotification) {
	data, err := json.Marshal(notification)
	if err != nil {
		log.Printf("failed to marshal notification: %v", err)
		return
	}
	_, _ = w.Write(data)
}
