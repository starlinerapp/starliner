package natscore

import (
	"context"
	"fmt"
	"time"

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

// SubscribeContext creates a short-lived subscription that is automatically
// torn down when ctx is canceled. Unlike Subscribe, it does not block the
// caller: the returned cancel function should be invoked to release the NATS
// subscription (defer is typical). The callback runs on NATS's internal
// delivery goroutine, so it must be non-blocking (e.g. a channel send).
func (s *Subscriber) SubscribeContext(
	ctx context.Context,
	subject Subject,
	identifier string,
	cb func([]byte),
) (func(), error) {
	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	sub, err := s.conn.Subscribe(uniqueSubject, func(msg *nats.Msg) {
		cb(msg.Data)
	})
	if err != nil {
		return nil, err
	}

	cancel := func() {
		_ = sub.Unsubscribe()
	}

	go func() {
		<-ctx.Done()
		cancel()
	}()

	return cancel, nil
}
