package notification

import (
	"sync"

	"starliner.app/internal/core/domain/value"
)

type userSubscription struct {
	hub    *UserNotificationHub
	userId int64
	ch     chan *value.ClusterNotification
	once   sync.Once
}

func (s *userSubscription) Notifications() <-chan *value.ClusterNotification {
	return s.ch
}

func (s *userSubscription) Close() {
	s.once.Do(func() { s.hub.unsubscribe(s.userId, s.ch) })
}

type environmentSubscription struct {
	hub  *EnvironmentNotificationHub
	key  environmentKey
	ch   chan *value.EnvironmentNotification
	once sync.Once
}

func (s *environmentSubscription) Notifications() <-chan *value.EnvironmentNotification {
	return s.ch
}

func (s *environmentSubscription) Close() {
	s.once.Do(func() { s.hub.unsubscribe(s.key, s.ch) })
}
