package pubsub

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"starliner.app/internal/builder/domain/port"
	natscore "starliner.app/internal/core/infrastructure/nats/core"
)

var Module = fx.Module(
	"pubsub",
	fx.Provide(
		natscore.Connect,
		func(conn *nats.Conn) port.LogPublisher {
			return NewPubsub(conn)
		},
	),
)
