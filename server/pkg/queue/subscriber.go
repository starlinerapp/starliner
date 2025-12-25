package queue

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

		cb(m)
		err := msg.Ack()
		if err != nil {
			return
		}
	},
		nats.Durable(durable),
		nats.ManualAck(),
		nats.AckWait(30*time.Second),
	)

	return err
}
