package sse

import (
	"encoding/json"
	"log"
	"sync"

	"starliner.app/internal/core/domain/value"
)

type UserNotificationHub struct {
	mu      sync.RWMutex
	clients map[int64]map[chan *value.ClusterNotification]struct{}
}

func NewUserNotificationHub() *UserNotificationHub {
	return &UserNotificationHub{
		clients: make(map[int64]map[chan *value.ClusterNotification]struct{}),
	}
}

func (h *UserNotificationHub) Subscribe(userId int64) chan *value.ClusterNotification {
	ch := make(chan *value.ClusterNotification, 16)
	h.mu.Lock()
	if h.clients[userId] == nil {
		h.clients[userId] = make(map[chan *value.ClusterNotification]struct{})
	}
	h.clients[userId][ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *UserNotificationHub) Unsubscribe(userId int64, ch chan *value.ClusterNotification) {
	h.mu.Lock()
	if subs, ok := h.clients[userId]; ok {
		delete(subs, ch)
		if len(subs) == 0 {
			delete(h.clients, userId)
		}
	}
	h.mu.Unlock()
	close(ch)
}

func (h *UserNotificationHub) Broadcast(userId int64, notification *value.ClusterNotification) {
	h.mu.RLock()
	subs := h.clients[userId]
	h.mu.RUnlock()

	for ch := range subs {
		select {
		case ch <- notification:
		default:
			log.Printf("dropping notification for slow client on user %d", userId)
		}
	}
}

func (h *UserNotificationHub) WriteNotification(w *Writer, notification *value.ClusterNotification) {
	data, err := json.Marshal(notification)
	if err != nil {
		log.Printf("failed to marshal user notification: %v", err)
		return
	}
	_, _ = w.Write(data)
}
