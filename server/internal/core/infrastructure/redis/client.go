package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"starliner.app/internal/core/domain/port"
)

type Client struct {
	client *redis.Client
}

func NewClient(client *redis.Client) *Client {
	return &Client{client: client}
}

func (c *Client) TryAcquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return c.client.SetNX(ctx, key, "1", ttl).Result()
}

func (c *Client) AppendToStream(ctx context.Context, name string, payload map[string][]byte) error {
	values := make(map[string]any, len(payload))

	for key, value := range payload {
		values[key] = value
	}

	return c.client.XAdd(ctx, &redis.XAddArgs{
		Stream: name,
		MaxLen: 10_000,
		Approx: true,
		Values: values,
	}).Err()
}

func (c *Client) ReadStream(ctx context.Context, name string, lastId string) ([]port.StreamEntry, error) {
	if lastId == "" {
		lastId = "0"
	}

	streams, err := c.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{name, lastId},
		Count:   100,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]port.StreamEntry, 0)

	for _, stream := range streams {
		for _, message := range stream.Messages {
			values := make(map[string][]byte, len(message.Values))

			for key, value := range message.Values {
				switch v := value.(type) {
				case string:
					values[key] = []byte(v)
				case []byte:
					values[key] = v
				}
			}

			entries = append(entries, port.StreamEntry{
				ID:     message.ID,
				Values: values,
			})
		}
	}

	return entries, nil
}
