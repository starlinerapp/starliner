package builder

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/internal/application"
	"starliner.app/internal/infrastructure/queue"
	"starliner.app/internal/infrastructure/queue/proto/v1"
)

type Consumer struct {
	buildSubscriber  *queue.Subscriber[*v1.Build]
	buildApplication *application.BuildApplication
}

func RegisterConsumer(lc fx.Lifecycle, o *Consumer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return o.Start()
		},
	})
}

func NewConsumer(
	buildSubscriber *queue.Subscriber[*v1.Build],
	buildApplication *application.BuildApplication,
) *Consumer {
	return &Consumer{
		buildSubscriber:  buildSubscriber,
		buildApplication: buildApplication,
	}
}

func (o *Consumer) Start() error {
	go func() {
		err := o.buildSubscriber.Subscribe(queue.BuildTriggered, "*", "buildTriggered", o.buildApplication.HandleBuildTriggered)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()
	return nil
}
