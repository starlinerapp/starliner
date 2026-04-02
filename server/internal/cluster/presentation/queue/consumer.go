package queue

import (
	"context"
	"go.uber.org/fx"
	"starliner.app/internal/cluster/application"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/util/concurrent"
)

type Consumer struct {
	deploymentApplication *application.DeploymentApplication
	imageApplication      *application.ImageApplication
	databaseApplication   *application.DatabaseApplication
	ingressApplication    *application.IngressApplication
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
	imageApplication *application.ImageApplication,
	databaseApplication *application.DatabaseApplication,
	ingressApplication *application.IngressApplication,
	queue port.Queue,
) *Consumer {
	return &Consumer{
		deploymentApplication: deploymentApplication,
		imageApplication:      imageApplication,
		databaseApplication:   databaseApplication,
		ingressApplication:    ingressApplication,
		queue:                 queue,
	}
}

func (c *Consumer) Start() error {
	go concurrent.WithRecovery(context.Background(), "SubscribeToDeployImage", func() error {
		return c.queue.SubscribeToDeployImage(c.imageApplication.HandleDeployImage)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToDeployDatabase", func() error {
		return c.queue.SubscribeToDeployDatabase(c.databaseApplication.HandleDeployDatabase)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToDeleteDeployment", func() error {
		return c.queue.SubscribeToDeleteDeployment(c.deploymentApplication.HandleDeleteDeployment)
	})

	go concurrent.WithRecovery(context.Background(), "SubscribeToDeployIngress", func() error {
		return c.queue.SubscribeToDeployIngress(c.ingressApplication.HandleDeployIngress)
	})

	return nil
}
