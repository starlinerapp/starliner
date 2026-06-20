package queue

import (
	"context"
	"log"

	"go.uber.org/fx"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/port"
	coreValue "starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/util/concurrent"
)

type Consumer struct {
	runnerApplication *application.RunnerApplication
	queue             port.Queue
}

func RegisterConsumer(lc fx.Lifecycle, c *Consumer) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return c.Start()
		},
	})
}

func NewConsumer(
	runnerApplication *application.RunnerApplication,
	queue port.Queue,
) *Consumer {
	return &Consumer{
		runnerApplication: runnerApplication,
		queue:             queue,
	}
}

func (c *Consumer) Start() error {
	go concurrent.WithRecovery(context.Background(), "SubscribeToRunnerStatusChanged", func() error {
		return c.queue.SubscribeToRunnerStatusChanged(func(status *coreValue.RunnerStatusChanged) {
			if err := c.runnerApplication.HandleRunnerStatusChanged(context.Background(), status); err != nil {
				log.Printf("failed to handle runner status changed: %v", err)
			}
		})
	})

	return nil
}
