package jetstream

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	js nats.JetStreamContext
}

func NewSubscriber(js nats.JetStreamContext) *Subscriber {
	return &Subscriber{js: js}
}

func (s *Subscriber) Subscribe(subject Subject, identifier string, durable string, cb func([]byte)) error {
	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	queueGroup := fmt.Sprintf("%s-workers", durable)
	sub, err := s.js.QueueSubscribe(uniqueSubject, queueGroup, func(msg *nats.Msg) {
		if err := msg.Ack(); err != nil {
			log.Printf("failed to ack message: %v", err)
			return
		}
		cb(msg.Data)
	},
		nats.Durable(durable),
		nats.ManualAck(),
		nats.AckWait(30*time.Second),
	)

	if err != nil {
		return err
	}

	for sub.IsValid() {
		time.Sleep(500 * time.Millisecond)
	}

	err = sub.Unsubscribe()
	if err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}
	return fmt.Errorf("subscription to %s lost", uniqueSubject)
}
