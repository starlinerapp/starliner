package port

import coreValue "starliner.app/internal/core/domain/value"

type UserNotificationPublisher interface {
	Broadcast(userId int64, notification *coreValue.ClusterNotification)
}

type EnvironmentNotificationPublisher interface {
	Broadcast(correlationId string, environmentId int64, notification *coreValue.EnvironmentNotification)
}

type UserNotificationSubscriber interface {
	Subscribe(userId int64) UserNotificationSubscription
}

type EnvironmentNotificationSubscriber interface {
	Subscribe(correlationId string, environmentId int64) EnvironmentNotificationSubscription
}

type UserNotificationSubscription interface {
	Notifications() <-chan *coreValue.ClusterNotification
	Close()
}

type EnvironmentNotificationSubscription interface {
	Notifications() <-chan *coreValue.EnvironmentNotification
	Close()
}
