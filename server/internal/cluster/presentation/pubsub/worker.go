package pubsub

import (
	"context"
	"go.uber.org/fx"
	"starliner.app/internal/cluster/application"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/util/concurrent"
)

type Worker struct {
	statusApplication *application.StatusApplication
	pubsub            port.Pubsub
}

func RegisterWorker(lc fx.Lifecycle, w *Worker) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return w.Start()
		},
	})
}

func NewWorker(statusApplication *application.StatusApplication, pubsub port.Pubsub) *Worker {
	return &Worker{
		statusApplication: statusApplication,
		pubsub:            pubsub,
	}
}

func (w *Worker) Start() error {
	go concurrent.WithRecovery(context.Background(), "SubscribeToDeploymentStatusRequest", func() error {
		return w.pubsub.SubscribeToDeploymentStatusRequest(w.statusApplication.HandleRequestDeploymentStatus)
	})

	return nil
}
