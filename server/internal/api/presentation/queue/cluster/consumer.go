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
	go concurrent.WithRecovery(context.Background(), "SubscribeToClusterProvisionedSuccess", func() error {
		return c.queue.SubscribeToClusterProvisionedSuccess(c.clusterApplication.HandleClusterProvisionedSuccess)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToClusterProvisionedFailure", func() error {
		return c.queue.SubscribeToClusterProvisionedFailure(c.clusterApplication.HandleClusterProvisionedFailure)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToClusterDeletedSuccess", func() error {
		return c.queue.SubscribeToClusterDeletedSuccess(c.clusterApplication.HandleClusterDeletedSuccess)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToClusterDeletedFailure", func() error {
		return c.queue.SubscribeToClusterDeletedFailure(c.clusterApplication.HandleClusterDeletedFailure)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToBuildSucceeded", func() error {
		return c.queue.SubscribeToBuildSucceeded(c.deploymentApplication.HandleBuildSucceeded)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToBuildFailed", func() error {
		return c.queue.SubscribeToBuildFailed(c.deploymentApplication.HandleBuildFailed)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToDatabaseDeployedSuccess", func() error {
		return c.queue.SubscribeToDatabaseDeployedSuccess(c.deploymentApplication.HandleDatabaseDeployedSuccess)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToDatabaseDeployedFailure", func() error {
		return c.queue.SubscribeToDatabaseDeployedFailure(c.deploymentApplication.HandleDatabaseDeployedFailure)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToImageDeployedSuccess", func() error {
		return c.queue.SubscribeToImageDeployedSuccess(c.deploymentApplication.HandleImageDeployedSuccess)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToImageDeployedFailure", func() error {
		return c.queue.SubscribeToImageDeployedFailure(c.deploymentApplication.HandleImageDeployedFailure)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToIngressDeployedSuccess", func() error {
		return c.queue.SubscribeToIngressDeployedSuccess(c.deploymentApplication.HandleIngressDeployedSuccess)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToIngressDeployedFailure", func() error {
		return c.queue.SubscribeToIngressDeployedFailure(c.deploymentApplication.HandleIngressDeployedFailure)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToDeploymentDeletedSuccess", func() error {
		return c.queue.SubscribeToDeploymentDeletedSuccess(c.deploymentApplication.HandleDeploymentDeletedSuccess)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToDeploymentDeletedFailure", func() error {
		return c.queue.SubscribeToDeploymentDeletedFailure(c.deploymentApplication.HandleDeploymentDeletedFailure)
	})


	go concurrent.WithRecovery(context.Background(), "SubscribeToDeploymentStatusLogsCompleted", func() error {
		return c.queue.SubscribeToDeploymentStatusLogsCompleted(c.deploymentApplication.HandleDeploymentStatusLogsCompleted)
	})

	return nil
}
