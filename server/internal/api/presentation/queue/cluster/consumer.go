package queue

import (
	"context"
	"go.uber.org/fx"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/core/util/concurrent"
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
	go concurrent.WithRecovery(context.Background(), "SubscribeToClusterCreated", func() error {
		return c.queue.SubscribeToClusterCreated(c.clusterApplication.HandleClusterCreated)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToClusterDeleted", func() error {
		return c.queue.SubscribeToClusterDeleted(c.clusterApplication.HandleClusterDeleted)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToBuildCompleted", func() error {
		return c.queue.SubscribeToBuildCompleted(c.deploymentApplication.HandleBuildCompleted)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToDatabaseDeploymentCreated", func() error {
		return c.queue.SubscribeToDatabaseDeploymentCreated(c.deploymentApplication.HandleDatabaseDeploymentCreated)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToDeploymentDeleted", func() error {
		return c.queue.SubscribeToDeploymentDeleted(c.deploymentApplication.HandleDeploymentDeleted)
	})

	return nil
}
