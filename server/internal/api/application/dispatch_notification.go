package application

import (
	"starliner.app/internal/api/domain/port"
	coreValue "starliner.app/internal/core/domain/value"
)

type DispatchNotification struct {
	delivery port.NotificationDelivery
}

func NewDispatchNotification(delivery port.NotificationDelivery) *DispatchNotification {
	return &DispatchNotification{delivery: delivery}
}

func (d *DispatchNotification) Execute(topic string, notification *coreValue.Notification) {
	d.delivery.Deliver(topic, notification)
}
