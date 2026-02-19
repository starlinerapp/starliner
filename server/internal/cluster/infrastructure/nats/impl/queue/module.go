package queue

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
)

const (
	Deployments jetstream.Stream = "deployments"
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
		return jetstream.EnsureStream(js, Deployments, []jetstream.Subject{
			DeployDatabase,
			DeleteDatabase,
			DatabaseDeleted,
			DeployIngress,
		})
	}),
)
