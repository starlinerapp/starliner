package nats

import (
	natsgo "github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"starliner.app/internal/domain/port"
	"starliner.app/internal/infrastructure/nats"
)

var Module = fx.Module(
	"queue",
	fx.Provide(
		nats.Connect,
		func(js natsgo.JetStreamContext) port.Queue {
			return NewQueue(js)
		},
	),
	fx.Invoke(func(js natsgo.JetStreamContext) error {
		return nats.EnsureStream(js, nats.Builds, []nats.Subject{nats.BuildTriggered})
	}),
	fx.Invoke(func(js natsgo.JetStreamContext) error {
		return nats.EnsureStream(js, nats.Clusters, []nats.Subject{nats.CreateCluster, nats.DeleteCluster})
	}),
	fx.Invoke(func(js natsgo.JetStreamContext) error {
		return nats.EnsureStream(js, nats.Projects, []nats.Subject{nats.CreateProject})
	}),
)
