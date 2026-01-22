package queue

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
	"starliner.app/internal/provisioner/domain/port"
)

const (
	Clusters jetstream.Stream = "clusters"
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
		return jetstream.EnsureStream(js, Clusters, []jetstream.Subject{CreateCluster, ClusterCreated, DeleteCluster, ClusterDeleted})
	}),
)
