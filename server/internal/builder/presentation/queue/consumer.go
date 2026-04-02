package queue

import (
	"context"
	"go.uber.org/fx"
	"starliner.app/internal/builder/application"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/util/concurrent"
)

type Consumer struct {
	queue            port.Queue
	buildApplication *application.BuildApplication
}

func RegisterConsumer(lc fx.Lifecycle, c *Consumer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return c.Start()
		},
	})
}

func NewConsumer(
	queue port.Queue,
	buildApplication *application.BuildApplication,
) *Consumer {
	return &Consumer{
		queue:            queue,
		buildApplication: buildApplication,
	}
}

func (c *Consumer) Start() error {
	go concurrent.WithRecovery(context.Background(), "SubscribeToBuildTriggered", func() error {
		return c.queue.SubscribeToBuildTriggered(c.buildApplication.HandleBuildTriggered)
	})
	return nil
}
