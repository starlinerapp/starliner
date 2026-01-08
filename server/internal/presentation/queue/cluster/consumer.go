package cluster

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/application"
	"starliner.app/internal/domain/port"
)

type Consumer struct {
	projectApplication *application.ProjectApplication
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
	projectApplication *application.ProjectApplication,
	queue port.Queue,
) *Consumer {
	return &Consumer{
		projectApplication: projectApplication,
		queue:              queue,
	}
}

func (o *Consumer) Start() error {
	go func() {
		err := o.queue.SubscribeToCreateProject(o.projectApplication.HandleCreateProject)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}
