package sse

import (
	"sync"

	"starliner.app/internal/core/domain/value"
)

type subscription struct {
	service *Service
	topics  []string
	ch      chan *value.Notification
	once    sync.Once
}

func (s *subscription) Notifications() <-chan *value.Notification {
	return s.ch
}

func (s *subscription) Close() {
	s.once.Do(func() { s.service.unsubscribe(s.topics, s.ch) })
}
