package application

import "starliner.app/internal/api/domain/port"

type NotificationApplication struct {
	userNotifications        port.UserNotificationSubscriber
	environmentNotifications port.EnvironmentNotificationSubscriber
}

func NewNotificationApplication(
	userNotifications port.UserNotificationSubscriber,
	environmentNotifications port.EnvironmentNotificationSubscriber,
) *NotificationApplication {
	return &NotificationApplication{
		userNotifications:        userNotifications,
		environmentNotifications: environmentNotifications,
	}
}

func (na *NotificationApplication) SubscribeUser(userId int64) port.UserNotificationSubscription {
	return na.userNotifications.Subscribe(userId)
}

func (na *NotificationApplication) SubscribeEnvironment(correlationId string, environmentId int64) port.EnvironmentNotificationSubscription {
	return na.environmentNotifications.Subscribe(correlationId, environmentId)
}
