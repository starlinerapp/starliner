package sse

import (
	"log"
	"sync"

	apiport "starliner.app/internal/api/domain/port"
	"starliner.app/internal/core/domain/value"
)

type Service struct {
	mu      sync.RWMutex
	clients map[string]map[chan *value.Notification]struct{}
}

func NewService() *Service {
	return &Service{
		clients: make(map[string]map[chan *value.Notification]struct{}),
	}
}

func (s *Service) Subscribe(topics ...string) apiport.NotificationSubscription {
	ch := make(chan *value.Notification, 16)
	s.mu.Lock()
	for _, t := range topics {
		if s.clients[t] == nil {
			s.clients[t] = make(map[chan *value.Notification]struct{})
		}
		s.clients[t][ch] = struct{}{}
	}
	s.mu.Unlock()
	return &subscription{service: s, topics: topics, ch: ch}
}

func (s *Service) unsubscribe(topics []string, ch chan *value.Notification) {
	s.mu.Lock()
	for _, t := range topics {
		if subs, ok := s.clients[t]; ok {
			delete(subs, ch)
			if len(subs) == 0 {
				delete(s.clients, t)
			}
		}
	}
	s.mu.Unlock()
}

func (s *Service) Deliver(topic string, notification *value.Notification) {
	s.mu.RLock()
	subs := s.clients[topic]
	targets := make([]chan *value.Notification, 0, len(subs))
	for ch := range subs {
		targets = append(targets, ch)
	}
	s.mu.RUnlock()

	for _, ch := range targets {
		select {
		case ch <- notification:
		default:
			log.Printf("dropping notification for slow client on topic %s", topic)
		}
	}
}
