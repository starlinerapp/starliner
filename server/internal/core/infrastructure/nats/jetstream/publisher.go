package jetstream

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	js nats.JetStreamContext
}

func NewPublisher(js nats.JetStreamContext) *Publisher {
	return &Publisher{js: js}
}

func (p *Publisher) Publish(subject Subject, identifier string, msg []byte) error {
	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	_, err := p.js.Publish(uniqueSubject, msg)
	return err
}
