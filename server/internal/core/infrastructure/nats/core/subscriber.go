package natscore

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
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
		if err := msg.Ack(); err != nil {
			log.Printf("failed to ack message: %v", err)
			return
		}

		cb(msg.Data)
	})

	return err
}
