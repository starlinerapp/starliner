package queue

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/cluster/application"
	"starliner.app/internal/cluster/domain/port"
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
	go func() {
		err := c.queue.SubscribeToDeployImage(c.imageApplication.HandleDeployImage)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	go func() {
		err := c.queue.SubscribeToDeployDatabase(c.databaseApplication.HandleDeployDatabase)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	go func() {
		err := c.queue.SubscribeToDeleteDeployment(c.deploymentApplication.HandleDeleteDeployment)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	go func() {
		err := c.queue.SubscribeToDeployIngress(c.ingressApplication.HandleDeployIngress)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}
