package cluster

import (
	"context"

	"go.uber.org/fx"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/core/util/concurrent"
)

type Subscriber struct {
	clusterApplication *application.ClusterApplication
	pubsub             port.Pubsub
}

func RegisterSubscriber(lc fx.Lifecycle, s *Subscriber) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return s.Start()
		},
	})
}

func NewSubscriber(
	clusterApplication *application.ClusterApplication,
	pubsub port.Pubsub,
) *Subscriber {
	return &Subscriber{
		clusterApplication: clusterApplication,
		pubsub:             pubsub,
	}
}

func (s *Subscriber) Start() error {
	go concurrent.WithRecovery(context.Background(), "SubscribeToReconcileClusterRequest", func() error {
		return s.pubsub.SubscribeToReconcileClusterRequest(s.clusterApplication.HandleReconcileClusterRequest)
	})

	return nil
}
