package queue

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/cluster/application"
	"starliner.app/internal/cluster/domain/port"
)

type Consumer struct {
	deploymentApplication *application.DeploymentApplication
	queue                 port.Queue
}

func RegisterConsumer(lc fx.Lifecycle, c *Consumer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return c.Start()
		},
	})
}

func NewConsumer(
	deploymentApplication *application.DeploymentApplication,
	queue port.Queue,
) *Consumer {
	return &Consumer{
		deploymentApplication: deploymentApplication,
		queue:                 queue,
	}
}

func (c *Consumer) Start() error {
	go func() {
		err := c.queue.SubscribeToDeployDatabase(c.deploymentApplication.HandleDeployDatabase)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	go func() {
		err := c.queue.SubscribeToDeleteDatabase(c.deploymentApplication.HandleDeleteDatabase)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}
