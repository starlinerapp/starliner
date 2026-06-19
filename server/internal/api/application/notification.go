package application

import (
	"fmt"

	"starliner.app/internal/api/domain/port"
)

func correlationTopic(correlationId string) string {
	return "corr:" + correlationId
}

func environmentTopic(correlationId string, environmentId int64) string {
	return fmt.Sprintf("corr:%s:env:%d", correlationId, environmentId)
}

type NotificationApplication struct {
	notifications port.NotificationSubscriber
}

func NewNotificationApplication(notifications port.NotificationSubscriber) *NotificationApplication {
	return &NotificationApplication{notifications: notifications}
}

func (na *NotificationApplication) SubscribeGlobal(correlationId string) port.NotificationSubscription {
	return na.notifications.Subscribe(correlationTopic(correlationId))
}

func (na *NotificationApplication) SubscribeEnvironment(correlationId string, environmentId int64) port.NotificationSubscription {
	return na.notifications.Subscribe(environmentTopic(correlationId, environmentId))
}
