package port

import "context"

type Subscription interface {
	Channel() <-chan []byte
	Close() error
}

type PubSub interface {
	Publish(ctx context.Context, channel string, payload []byte) error
	Subscribe(ctx context.Context, channel string) (Subscription, error)
}
