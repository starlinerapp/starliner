package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Publisher[T proto.Message] struct {
	js nats.JetStreamContext
}

func NewPublisher[T proto.Message](js nats.JetStreamContext) *Publisher[T] {
	return &Publisher[T]{js: js}
}

func (p *Publisher[T]) Publish(subject Subject, identifier string, msg T) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	_, err = p.js.Publish(uniqueSubject, data)
	return err
}
