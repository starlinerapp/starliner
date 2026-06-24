package notification

import (
	"context"

	"go.uber.org/fx"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/core/util/concurrent"
)

type Consumer struct {
	source     port.NotificationConsumer
	dispatcher *application.DispatchNotification
}

func NewConsumer(source port.NotificationConsumer, dispatcher *application.DispatchNotification) *Consumer {
	return &Consumer{source: source, dispatcher: dispatcher}
}

func RegisterConsumer(lc fx.Lifecycle, c *Consumer) {
	ctx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go concurrent.WithRecovery(ctx, "ConsumeNotifications", func() error {
				return c.source.Consume(ctx, c.dispatcher.Execute)
			})
			return nil
		},
		OnStop: func(context.Context) error {
			cancel()
			return nil
		},
	})
}
