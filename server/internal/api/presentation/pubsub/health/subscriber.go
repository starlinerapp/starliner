package health

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/port"
)

type Subscriber struct {
	deploymentApplication *application.DeploymentApplication
	pubsub                port.Pubsub
}

func RegisterSubscriber(lc fx.Lifecycle, s *Subscriber) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return s.Start()
		},
	})
}

func NewSubscriber(
	deploymentApplication *application.DeploymentApplication,
	pubsub port.Pubsub,
) *Subscriber {
	return &Subscriber{
		deploymentApplication: deploymentApplication,
		pubsub:                pubsub,
	}
}

func (s *Subscriber) Start() error {
	go func() {
		err := s.pubsub.SubscribeToDeploymentStatusResponse(s.deploymentApplication.HandleDeploymentStatusResponse)
		if err != nil {
			log.Fatalf("failed to subscribe to pubsub: %v", err)
		}
	}()
	return nil
}
