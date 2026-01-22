package queue

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/port"
)

type Consumer struct {
	deploymentApplication *application.DeploymentApplication
	clusterApplication    *application.ClusterApplication
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
	clusterApplication *application.ClusterApplication,
	queue port.Queue,
) *Consumer {
	return &Consumer{
		deploymentApplication: deploymentApplication,
		clusterApplication:    clusterApplication,
		queue:                 queue,
	}
}

func (c *Consumer) Start() error {
	go func() {
		err := c.queue.SubscribeToClusterCreated(c.clusterApplication.HandleClusterCreated)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()

	go func() {
		err := c.queue.SubscribeToClusterDeleted(c.clusterApplication.HandleClusterDeleted)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()

	go func() {
		err := c.queue.SubscribeToDatabaseDeleted(c.deploymentApplication.HandleDatabaseDeleted)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}
