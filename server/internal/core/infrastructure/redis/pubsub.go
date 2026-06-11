package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"starliner.app/internal/core/domain/port"
)

type PubSub struct {
	client *redis.Client
}

func NewPubSub(client *redis.Client) *PubSub {
	return &PubSub{client: client}
}

func (p *PubSub) Publish(ctx context.Context, channel string, payload []byte) error {
	return p.client.Publish(ctx, channel, payload).Err()
}

func (p *PubSub) Subscribe(ctx context.Context, channel string) (port.Subscription, error) {
	sub := p.client.Subscribe(ctx, channel)

	if _, err := sub.Receive(ctx); err != nil {
		_ = sub.Close()
		return nil, err
	}

	out := make(chan []byte, 64)

	go func() {
		defer close(out)
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				_ = sub.Close()
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				select {
				case out <- []byte(msg.Payload):
				default:
					log.Printf("redis pubsub: dropping message on channel %q (slow consumer)", channel)
				}
			}
		}
	}()

	return &redisSubscription{sub: sub, out: out}, nil
}

type redisSubscription struct {
	sub *redis.PubSub
	out chan []byte
}

func (s *redisSubscription) Channel() <-chan []byte { return s.out }
func (s *redisSubscription) Close() error           { return s.sub.Close() }
