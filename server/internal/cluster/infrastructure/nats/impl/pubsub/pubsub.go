package pubsub

import (
	"github.com/nats-io/nats.go"
	"starliner.app/internal/cluster/domain/port"
)

type Pubsub struct{}

func NewPubsub(conn *nats.Conn) port.Pubsub {
	return &Pubsub{}
}
