package cluster

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/application"
	"starliner.app/internal/domain/port"
)

type Consumer struct {
	environmentApplication *application.EnvironmentApplication
	queue                  port.Queue
}

func RegisterConsumer(lc fx.Lifecycle, o *Consumer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return o.Start()
		},
	})
}

func NewConsumer(
	environmentApplication *application.EnvironmentApplication,
	queue port.Queue,
) *Consumer {
	return &Consumer{
		environmentApplication: environmentApplication,
		queue:                  queue,
	}
}

func (o *Consumer) Start() error {
	go func() {
		err := o.queue.SubscribeToDeployDatabase(o.environmentApplication.HandleDeployDatabase)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}
