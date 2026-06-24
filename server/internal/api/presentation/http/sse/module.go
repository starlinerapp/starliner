package sse

import (
	"go.uber.org/fx"
	"starliner.app/internal/api/domain/port"
)

var Module = fx.Module(
	"sse",
	fx.Provide(
		NewService,
		func(s *Service) port.NotificationSubscriber { return s },
		func(s *Service) port.NotificationDelivery { return s },
	),
)
