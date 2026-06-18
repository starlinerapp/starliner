package notification

import (
	"go.uber.org/fx"
	"starliner.app/internal/api/domain/port"
)

var Module = fx.Module(
	"notification",
	fx.Provide(
		NewUserNotificationHub,
		NewEnvironmentNotificationHub,
		func(h *UserNotificationHub) port.UserNotificationPublisher { return h },
		func(h *UserNotificationHub) port.UserNotificationSubscriber { return h },
		func(h *EnvironmentNotificationHub) port.EnvironmentNotificationPublisher { return h },
		func(h *EnvironmentNotificationHub) port.EnvironmentNotificationSubscriber { return h },
	),
)
