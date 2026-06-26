package port

import (
	"context"
	"time"
)

type Lease interface {
	TryLease(ctx context.Context, key string, ttl time.Duration) (bool, error)
}
