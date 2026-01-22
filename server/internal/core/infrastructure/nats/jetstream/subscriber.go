package jetstream

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type Subscriber struct {
	js nats.JetStreamContext
}

func NewSubscriber(js nats.JetStreamContext) *Subscriber {
	return &Subscriber{js: js}
}

func (s *Subscriber) Subscribe(subject Subject, identifier string, durable string, cb func([]byte)) error {
	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	_, err := s.js.Subscribe(uniqueSubject, func(msg *nats.Msg) {
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

	return err
}
