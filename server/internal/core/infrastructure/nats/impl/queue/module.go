package queue

import (
	natsgo "github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
)

const (
	Builds      jetstream.Stream = "builds"
	Clusters    jetstream.Stream = "clusters"
	Deployments jetstream.Stream = "deployments"
)

var Module = fx.Module(
	"queue",
	fx.Provide(
		jetstream.Connect,
		func(js natsgo.JetStreamContext) port.Queue {
			return NewQueue(js)
		},
	),
	fx.Invoke(func(js natsgo.JetStreamContext) error {
		return jetstream.EnsureStream(js, Builds, []jetstream.Subject{BuildTriggered})
	}),
	fx.Invoke(func(js natsgo.JetStreamContext) error {
		return jetstream.EnsureStream(js, Clusters, []jetstream.Subject{CreateCluster, ClusterCreated, DeleteCluster, ClusterDeleted})
	}),
	fx.Invoke(func(js natsgo.JetStreamContext) error {
		return jetstream.EnsureStream(js, Deployments, []jetstream.Subject{DeployDatabase, DeleteDatabase, DatabaseDeleted})
	}),
)
