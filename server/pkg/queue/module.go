package queue

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	v1 "starliner.app/pkg/proto/v1"
)

var Module = fx.Module(
	"queue",
	fx.Provide(
		Connect,
		func(js nats.JetStreamContext) *Publisher[*v1.Project] {
			return NewPublisher[*v1.Project](js)
		},
		func(js nats.JetStreamContext) *Subscriber[*v1.Project] {
			return NewSubscriber[*v1.Project](js)
		},
	),
	fx.Invoke(func(js nats.JetStreamContext) error {
		return EnsureStream(js, Projects, []Subject{ProjectCreated})
	}),
)
