package queue

import (
	"context"
	"log"

	"go.uber.org/fx"
	"starliner.app/internal/builder/application"
	"starliner.app/internal/builder/domain/port"
	coreValue "starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/util/concurrent"
)

type Consumer struct {
	queue                port.Queue
	buildApplication     *application.BuildApplication
	heartbeatApplication *application.HeartbeatApplication
	runnerApplication    *application.RunnerApplication
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
	heartbeatApplication *application.HeartbeatApplication,
	runnerApplication *application.RunnerApplication,
) *Consumer {
	return &Consumer{
		queue:                queue,
		buildApplication:     buildApplication,
		heartbeatApplication: heartbeatApplication,
		runnerApplication:    runnerApplication,
	}
}

func (c *Consumer) Start() error {
	go concurrent.WithRecovery(context.Background(), "SubscribeToBuildTriggered", func() error {
		return c.queue.SubscribeToBuildTriggered(c.buildApplication.HandleBuildTriggered)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToRunnerDeleted", func() error {
		return c.queue.SubscribeToRunnerDeleted(func(runner *coreValue.RunnerDeleted) {
			if err := c.runnerApplication.HandleRunnerDeleted(context.Background(), runner.RunnerId); err != nil {
				log.Printf("failed to handle runner deleted: %v", err)
			}
		})
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToRunnerJobResults", func() error {
		return c.queue.SubscribeToRunnerJobResults(c.buildApplication.HandleRunnerJobResult)
	})

	return nil
}
