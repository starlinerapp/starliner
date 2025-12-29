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
		func(js nats.JetStreamContext) *Publisher[*v1.Build] {
			return NewPublisher[*v1.Build](js)
		},
		func(js nats.JetStreamContext) *Subscriber[*v1.Build] {
			return NewSubscriber[*v1.Build](js)
		},
		func(js nats.JetStreamContext) *Publisher[*v1.Cluster] {
			return NewPublisher[*v1.Cluster](js)
		},
		func(js nats.JetStreamContext) *Subscriber[*v1.Cluster] {
			return NewSubscriber[*v1.Cluster](js)
		},
	),
	fx.Invoke(func(js nats.JetStreamContext) error {
		return EnsureStream(js, Builds, []Subject{BuildTriggered})
	}),
	fx.Invoke(func(js nats.JetStreamContext) error {
		return EnsureStream(js, Clusters, []Subject{CreateCluster, DeleteCluster})
	}),
)
