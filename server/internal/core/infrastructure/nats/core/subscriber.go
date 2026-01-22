package natscore

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	conn *nats.Conn
}

func NewSubscriber(conn *nats.Conn) *Subscriber {
	return &Subscriber{conn: conn}
}

func (s *Subscriber) Subscribe(subject Subject, identifier string, cb func([]byte)) error {
	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	_, err := s.conn.Subscribe(uniqueSubject, func(msg *nats.Msg) {
		cb(msg.Data)
	})

	return err
}
