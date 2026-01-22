package queue

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
)

const (
	BuildTriggered jetstream.Subject = "build.triggered"
)

type Queue struct {
	subscriber *jetstream.Subscriber
}

func NewQueue(js nats.JetStreamContext) port.Queue {
	return &Queue{
		subscriber: jetstream.NewSubscriber(js),
	}
}

func (q *Queue) SubscribeToBuildTriggered(handler func(build *value.TriggerBuild)) error {
	return q.subscriber.Subscribe(BuildTriggered, "*", "buildTriggered", func(msg []byte) {
		var b value.TriggerBuild
		if err := json.Unmarshal(msg, &b); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&b)
	})
}
