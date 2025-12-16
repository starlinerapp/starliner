package queue

import (
	"github.com/nats-io/nats.go"
	"log"
)

func Connect() (nats.JetStreamContext, error) {
	nc, err := nats.Connect("http://nats:4222")
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
