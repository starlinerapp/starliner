package jetstream

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	js nats.JetStreamContext
}

func NewSubscriber(js nats.JetStreamContext) *Subscriber {
	return &Subscriber{js: js}
}

func (s *Subscriber) Subscribe(subject Subject, identifier string, durable string, cb func([]byte)) error {
	uniqueSubject := fmt.Sprintf("%s.%s", subject, identifier)
	queueGroup := fmt.Sprintf("%s-workers", durable)
	sub, err := s.js.QueueSubscribe(uniqueSubject, queueGroup, func(msg *nats.Msg) {
		go func() {
			done := make(chan struct{})
			go func() {
				ticker := time.NewTicker(15 * time.Second)
				defer ticker.Stop()
				if err := msg.InProgress(); err != nil {
					log.Printf("failed to set message in progress: %v", err)
				}
				for {
					select {
					case <-done:
						return
					case <-ticker.C:
						if err := msg.InProgress(); err != nil {
							log.Printf("failed to set message in progress: %v", err)
							return
						}
					}
				}
			}()

			cb(msg.Data)
			close(done)

			if err := msg.Ack(); err != nil {
				log.Printf("failed to ack message: %v", err)
			}
		}()
	},
		nats.Durable(durable),
		nats.ManualAck(),
		nats.AckWait(30*time.Second),
	)

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
