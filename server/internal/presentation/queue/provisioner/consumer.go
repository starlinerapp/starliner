package provisioner

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/application"
	"starliner.app/internal/infrastructure/nats"
	v1 "starliner.app/internal/infrastructure/nats/proto/v1"
)

type Consumer struct {
	clusterApplication *application.ClusterApplication
	clusterSubscriber  *nats.Subscriber[*v1.Cluster]
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
	clusterSubscriber *nats.Subscriber[*v1.Cluster],
) *Consumer {
	return &Consumer{
		clusterApplication: clusterApplication,
		clusterSubscriber:  clusterSubscriber,
	}
}

func (o *Consumer) Start() error {
	go func() {
		err := o.clusterSubscriber.Subscribe(nats.CreateCluster, "*", "createCluster", o.clusterApplication.HandleCreateCluster)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()

	go func() {
		err := o.clusterSubscriber.Subscribe(nats.DeleteCluster, "*", "deleteCluster", o.clusterApplication.HandleDeleteCluster)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}
