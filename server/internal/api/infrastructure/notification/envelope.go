package notification

import "starliner.app/internal/core/domain/value"

const notificationChannel = "notifications"

type envelope struct {
	Topic        string              `json:"topic"`
	Notification *value.Notification `json:"notification"`
}
