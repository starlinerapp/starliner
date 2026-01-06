package provisioner

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/application"
	"starliner.app/internal/domain/port"
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

func (o *Consumer) Start() error {
	go func() {
		err := o.queue.SubscribeToCreateCluster(o.clusterApplication.HandleCreateCluster)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()

	go func() {
		err := o.queue.SubscribeToDeleteCluster(o.clusterApplication.HandleDeleteCluster)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}
