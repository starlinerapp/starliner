package pubsub

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/cluster/application"
	"starliner.app/internal/cluster/domain/port"
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
	go func() {
		err := w.pubsub.SubscribeToDeploymentStatusRequest(w.statusApplication.HandleRequestDeploymentStatus)
		if err != nil {
			log.Fatalf("failed to subscribe to pubsub: %v", err)
		}
	}()
	return nil
}
