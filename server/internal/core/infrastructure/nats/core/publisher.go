package natscore

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type Subject string

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(conn *nats.Conn) *Publisher {
	return &Publisher{conn: conn}
}

func (p *Publisher) Publish(subject Subject, identifier string, msg []byte) error {
	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	return p.conn.Publish(uniqueSubject, msg)
}
