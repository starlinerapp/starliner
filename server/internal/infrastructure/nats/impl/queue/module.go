package queue

import (
	natsgo "github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"starliner.app/internal/domain/port"
	"starliner.app/internal/infrastructure/nats"
)

const (
	Builds   nats.Stream = "builds"
	Clusters nats.Stream = "clusters"
	Projects nats.Stream = "projects"
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
		return nats.EnsureStream(js, Builds, []nats.Subject{BuildTriggered})
	}),
	fx.Invoke(func(js natsgo.JetStreamContext) error {
		return nats.EnsureStream(js, Clusters, []nats.Subject{CreateCluster, DeleteCluster})
	}),
	fx.Invoke(func(js natsgo.JetStreamContext) error {
		return nats.EnsureStream(js, Projects, []nats.Subject{CreateProject})
	}),
)
