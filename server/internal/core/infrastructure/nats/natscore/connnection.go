package natscore

import (
	"github.com/nats-io/nats.go"
	"log"
	"starliner.app/internal/core/conf"
)

func Connect(cfg conf.NatsConfig) (*nats.Conn, error) {
	nc, err := nats.Connect(cfg.GetNatsUrl())
	if err != nil {
		log.Printf("failed to connect to NATS: %v", err)
		return nil, err
	}

	return nc, nil
}
