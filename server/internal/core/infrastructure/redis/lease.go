package redis

import (
	"context"
	"time"

	"starliner.app/internal/core/domain/port"
)

func (c *Client) TryLease(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return c.client.SetNX(ctx, key, "1", ttl).Result()
}

var _ port.Lease = (*Client)(nil)
