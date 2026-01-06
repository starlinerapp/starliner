package cluster

import (
	"context"
	_ "embed"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/application"
	"starliner.app/internal/infrastructure/nats"
	"starliner.app/internal/infrastructure/nats/proto/v1"
)

type Consumer struct {
	projectApplication *application.ProjectApplication
	projectSubscriber  *nats.Subscriber[*v1.Project]
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
	projectSubscriber *nats.Subscriber[*v1.Project],
) *Consumer {
	return &Consumer{
		projectApplication: projectApplication,
		projectSubscriber:  projectSubscriber,
	}
}

func (o *Consumer) Start() error {
	go func() {
		err := o.projectSubscriber.Subscribe(nats.CreateProject, "*", "createProject", o.projectApplication.HandleCreateProject)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}
