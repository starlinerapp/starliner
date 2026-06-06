package port

import "context"

// Subscription represents a live subscription to a Pub/Sub channel.
// Channel returns a stream of payloads. Close terminates the subscription
// and releases the underlying connection.
type Subscription interface {
	Channel() <-chan []byte
	Close() error
}

// PubSub is a fire-and-forget fanout primitive used to broadcast messages
// across all processes that subscribe to the same channel. Messages are
// delivered to every subscriber (including the publisher), with at-most-once
// semantics. It is intended for ephemeral signals such as live notifications,
// not for durable message delivery — use KVStore (streams) for that.
type PubSub interface {
	Publish(ctx context.Context, channel string, payload []byte) error
	Subscribe(ctx context.Context, channel string) (Subscription, error)
}

