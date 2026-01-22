package pubsub

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/infrastructure/nats/natscore"
)

var Module = fx.Module(
	"pubsub",
	fx.Provide(
		natscore.Connect,
		func(p *nats.Conn) port.Pubsub {
			return NewPubsub(p)
		},
	),
)
