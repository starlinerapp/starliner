package port

import (
	"context"

	coreValue "starliner.app/internal/core/domain/value"
)

type NotificationPublisher interface {
	Publish(topic string, notification *coreValue.Notification)
}

type NotificationConsumer interface {
	Consume(ctx context.Context, handle func(topic string, notification *coreValue.Notification)) error
}

type NotificationSubscriber interface {
	Subscribe(topics ...string) NotificationSubscription
}

type NotificationDelivery interface {
	Deliver(topic string, notification *coreValue.Notification)
}

type NotificationSubscription interface {
	Notifications() <-chan *coreValue.Notification
	Close()
}
