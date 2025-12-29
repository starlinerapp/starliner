package cluster

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/pkg/config"
	v1 "starliner.app/pkg/proto/v1"
	"starliner.app/pkg/queue"
)

type Orchestrator struct {
	cfg               *config.Config
	clusterSubscriber *queue.Subscriber[*v1.Cluster]
}

func RegisterOrchestrator(lc fx.Lifecycle, o *Orchestrator) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return o.Start()
		},
	})
}

func NewOrchestrator(
	cfg *config.Config,
	clusterSubscriber *queue.Subscriber[*v1.Cluster],
) *Orchestrator {
	return &Orchestrator{
		cfg:               cfg,
		clusterSubscriber: clusterSubscriber,
	}
}

func (o *Orchestrator) Start() error {
	go func() {
		err := o.clusterSubscriber.Subscribe(queue.CreateCluster, "*", "createCluster", o.handleCreateCluster)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()

	return nil
}

func (o *Orchestrator) handleCreateCluster(cluster *v1.Cluster) {
	log.Printf("create cluster: %s", cluster.Name)
}
