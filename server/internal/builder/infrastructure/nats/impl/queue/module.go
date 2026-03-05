package queue

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
)

const (
	Builds jetstream.Stream = "builds"
)

var Module = fx.Module(
	"queue",
	fx.Provide(
		jetstream.Connect,
		func(js nats.JetStreamContext) port.Queue {
			return NewQueue(js)
		},
	),
	fx.Invoke(func(js nats.JetStreamContext) error {
		return jetstream.EnsureStream(js, Builds, []jetstream.Subject{BuildTriggered, BuildCompleted})
	}),
)
