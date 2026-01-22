package natscore

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"log"
	"reflect"
)

type Subscriber[T proto.Message] struct {
	conn *nats.Conn
}

func NewSubscriber[T proto.Message](conn *nats.Conn) *Subscriber[T] {
	return &Subscriber[T]{conn: conn}
}

func (s *Subscriber[T]) Subscribe(subject Subject, identifier string, cb func(T)) error {
	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	_, err := s.conn.Subscribe(uniqueSubject, func(msg *nats.Msg) {
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
	})

	return err
}
