package notification

import (
	"go.uber.org/fx"
	"starliner.app/internal/api/domain/port"
)

var Module = fx.Module(
	"notification",
	fx.Provide(
		NewPublisher,
		func(p *Publisher) port.NotificationPublisher { return p },
		NewSource,
		func(s *Source) port.NotificationConsumer { return s },
	),
)
