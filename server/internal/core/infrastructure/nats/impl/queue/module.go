package queue

import (
	natsgo "github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/infrastructure/nats"
)

const (
	Builds      nats.Stream = "builds"
	Clusters    nats.Stream = "clusters"
	Deployments nats.Stream = "deployments"
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
		return nats.EnsureStream(js, Clusters, []nats.Subject{CreateCluster, ClusterCreated, DeleteCluster, ClusterDeleted})
	}),
	fx.Invoke(func(js natsgo.JetStreamContext) error {
		return nats.EnsureStream(js, Deployments, []nats.Subject{DeployDatabase, DeleteDatabase})
	}),
)
