package notification

import (
	"context"
	"encoding/json"
	"log"

	"starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

type Publisher struct {
	pubsub port.PubSub
}

func NewPublisher(ps port.PubSub) *Publisher {
	return &Publisher{pubsub: ps}
}

func (p *Publisher) Publish(topic string, notification *value.Notification) {
	payload, err := json.Marshal(envelope{Topic: topic, Notification: notification})
	if err != nil {
		log.Printf("failed to marshal notification envelope: %v", err)
		return
	}

	if err := p.pubsub.Publish(context.Background(), notificationChannel, payload); err != nil {
		log.Printf("failed to publish notification: %v", err)
	}
}
