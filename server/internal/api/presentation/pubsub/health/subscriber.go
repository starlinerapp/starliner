package health

import (
	"context"
	"go.uber.org/fx"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/core/util/concurrent"
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
	go concurrent.WithRecovery(context.Background(), "SubscribeToDeploymentStatusResponse", func() error {
		return s.pubsub.SubscribeToDeploymentStatusResponse(s.deploymentApplication.HandleDeploymentStatusResponse)
	})

	return nil
}
