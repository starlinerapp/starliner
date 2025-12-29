package cluster

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"log"
	"starliner.app/pkg/config"
	v1 "starliner.app/pkg/proto/v1"
	"starliner.app/pkg/queue"
	interfaces "starliner.app/pkg/repository/interface"
)

type Orchestrator struct {
	cfg               *config.Config
	clusterRepository interfaces.ClusterRepository
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
	clusterRepository interfaces.ClusterRepository,
	clusterSubscriber *queue.Subscriber[*v1.Cluster],
) *Orchestrator {
	return &Orchestrator{
		cfg:               cfg,
		clusterRepository: clusterRepository,
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

	go func() {
		err := o.clusterSubscriber.Subscribe(queue.DeleteCluster, "*", "deleteCluster", o.handleDeleteCluster)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}

func (o *Orchestrator) handleCreateCluster(c *v1.Cluster) {
	fmt.Printf("create cluster: %s\n", c.Name)
}

func (o *Orchestrator) handleDeleteCluster(c *v1.Cluster) {
	ctx := context.Background()
	if c.Id == nil {
		fmt.Printf("cannot delete cluster: missing cluster id")
		return
	}
	err := o.clusterRepository.DeleteCluster(ctx, *c.Id)
	if err != nil {
		fmt.Printf("failed to delete cluster from database: %v", err)
	}
}
