package builder

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/application"
	"starliner.app/internal/domain/port"
)

type Consumer struct {
	queue            port.Queue
	buildApplication *application.BuildApplication
}

func RegisterConsumer(lc fx.Lifecycle, o *Consumer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return o.Start()
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
	go func() {
		err := c.queue.SubscribeToBuildTriggered(c.buildApplication.HandleBuildTriggered)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}
