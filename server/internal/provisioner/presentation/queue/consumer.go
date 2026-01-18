package queue

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/core/domain/port"
	"starliner.app/internal/provisioner/application"
)

type Consumer struct {
	clusterApplication *application.ClusterApplication
	queue              port.Queue
}

func RegisterConsumer(lc fx.Lifecycle, o *Consumer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return o.Start()
		},
	})
}

func NewConsumer(
	clusterApplication *application.ClusterApplication,
	queue port.Queue,
) *Consumer {
	return &Consumer{
		clusterApplication: clusterApplication,
		queue:              queue,
	}
}

func (c *Consumer) Start() error {
	go func() {
		err := c.queue.SubscribeToCreateCluster(c.clusterApplication.HandleProvisionCluster)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()

	go func() {
		err := c.queue.SubscribeToDeleteCluster(c.clusterApplication.HandleDeleteCluster)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}
