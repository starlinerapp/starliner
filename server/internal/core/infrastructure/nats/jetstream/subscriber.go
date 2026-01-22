package jetstream

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"log"
	"reflect"
	"time"
)

type Subscriber[T proto.Message] struct {
	js nats.JetStreamContext
}

func NewSubscriber[T proto.Message](js nats.JetStreamContext) *Subscriber[T] {
	return &Subscriber[T]{js: js}
}

func (s *Subscriber[T]) Subscribe(subject Subject, identifier string, durable string, cb func(T)) error {
	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	_, err := s.js.Subscribe(uniqueSubject, func(msg *nats.Msg) {
		var t T
		m := reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)

		if err := proto.Unmarshal(msg.Data, m); err != nil {
			log.Printf("failed to decode message: %v", err)
			err := msg.Nak()
			if err != nil {
				return
			}
			return
		}
		if err := msg.Ack(); err != nil {
			log.Printf("failed to ack message: %v", err)
			return
		}

		cb(m)
	},
		nats.Durable(durable),
		nats.ManualAck(),
		nats.AckWait(30*time.Second),
	)

	return err
}
