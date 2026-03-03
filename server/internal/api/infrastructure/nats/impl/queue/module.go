package queue

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"starliner.app/internal/api/domain/port"
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
		func(js nats.JetStreamContext) port.Queue {
			return NewQueue(js)
		},
	),
	fx.Invoke(func(js nats.JetStreamContext) error {
		return jetstream.EnsureStream(js, Builds, []jetstream.Subject{BuildTriggered})
	}),
	fx.Invoke(func(js nats.JetStreamContext) error {
		return jetstream.EnsureStream(js, Clusters, []jetstream.Subject{CreateCluster, ClusterCreated, DeleteCluster, ClusterDeleted})
	}),
	fx.Invoke(func(js nats.JetStreamContext) error {
		return jetstream.EnsureStream(js, Deployments, []jetstream.Subject{DeployImage, DeployDatabase, DatabaseDeployed, DeleteDeployment, DeploymentDeleted})
	}),
)
