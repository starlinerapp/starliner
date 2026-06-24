package notification

import (
	"context"
	"encoding/json"
	"log"

	"starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

type Source struct {
	pubsub port.PubSub
}

func NewSource(ps port.PubSub) *Source {
	return &Source{pubsub: ps}
}

func (s *Source) Consume(ctx context.Context, handle func(topic string, notification *value.Notification)) error {
	sub, err := s.pubsub.Subscribe(ctx, notificationChannel)
	if err != nil {
		return err
	}
	defer func() {
		if err := sub.Close(); err != nil {
			log.Printf("failed to close notification subscription: %v", err)
		}
	}()

	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case raw, ok := <-ch:
			if !ok {
				return nil
			}
			var env envelope
			if err := json.Unmarshal(raw, &env); err != nil {
				log.Printf("failed to unmarshal notification envelope: %v", err)
				continue
			}
			handle(env.Topic, env.Notification)
		}
	}
}
