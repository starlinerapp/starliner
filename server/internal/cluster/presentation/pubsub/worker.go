package pubsub

import (
	"context"
	"go.uber.org/fx"
	"starliner.app/internal/core/domain/port"
)

type Worker struct {
	pubsub port.Pubsub
}

func RegisterWorker(lc fx.Lifecycle, w *Worker) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return w.Start()
		},
	})
}

func NewWorker(pubsub port.Pubsub) *Worker {
	return &Worker{
		pubsub: pubsub,
	}
}

func (w *Worker) Start() error {
	go func() {
	}()
	return nil
}
