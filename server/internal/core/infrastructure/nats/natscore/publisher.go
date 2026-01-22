package natscore

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Subject string

type Publisher[T proto.Message] struct {
	conn *nats.Conn
}

func NewPublisher[T proto.Message](conn *nats.Conn) *Publisher[T] {
	return &Publisher[T]{conn: conn}
}

func (p *Publisher[T]) Publish(subject Subject, identifier string, msg T) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	return p.conn.Publish(uniqueSubject, data)
}
