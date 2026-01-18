package nats

import (
	"github.com/nats-io/nats.go"
	"log"
	"starliner.app/internal/core/conf"
)

func Connect(cfg *conf.Config) (nats.JetStreamContext, error) {
	nc, err := nats.Connect(cfg.NatsUrl)
	if err != nil {
		log.Printf("failed to connect to NATS: %v", err)
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, err
	}

	return js, nil
}
