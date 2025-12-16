package queue

import (
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type Subscriber[T proto.Message] struct {
	js nats.JetStreamContext
}

func NewSubscriber[T proto.Message](js nats.JetStreamContext) *Subscriber[T] {
	return &Subscriber[T]{js: js}
}

func (s *Subscriber[T]) Subscribe(subject string, durable string, cb func(T)) error {
	_, err := s.js.Subscribe(subject, func(msg *nats.Msg) {
		var m T
		if err := proto.Unmarshal(msg.Data, m); err != nil {
			log.Printf("failed to decode message: %v", err)
			msg.Nak()
			return
		}

		cb(m)
		msg.Ack()
	},
		nats.Durable(durable),
		nats.ManualAck(),
		nats.AckWait(30*time.Second),
	)

	return err
}
