package natscore

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

type Subscriber struct {
	conn *nats.Conn
}

func NewSubscriber(conn *nats.Conn) *Subscriber {
	return &Subscriber{conn: conn}
}

func (s *Subscriber) Subscribe(subject Subject, identifier string, cb func([]byte)) error {
	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	sub, err := s.conn.Subscribe(uniqueSubject, func(msg *nats.Msg) {
		cb(msg.Data)
	})

	if err != nil {
		return err
	}

	for sub.IsValid() {
		time.Sleep(500 * time.Millisecond)
	}

	err = sub.Unsubscribe()
	if err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}
	return fmt.Errorf("subscription to %s lost", uniqueSubject)
}
