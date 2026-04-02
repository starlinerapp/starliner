package queue

import (
	"context"
	"go.uber.org/fx"
	"starliner.app/internal/core/util/concurrent"
	"starliner.app/internal/provisioner/application"
	"starliner.app/internal/provisioner/domain/port"
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
	go concurrent.WithRecovery(context.Background(), "SubscribeToCreateCluster", func() error {
		return c.queue.SubscribeToCreateCluster(c.clusterApplication.HandleProvisionCluster)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToDeleteCluster", func() error {
		return c.queue.SubscribeToDeleteCluster(c.clusterApplication.HandleDeleteCluster)
	})

	return nil
}
