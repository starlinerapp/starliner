package port

import (
	"context"
	"time"
)

// AcquireLimiter gates an operation with a TTL-bound key across processes.
type AcquireLimiter interface {
	TryAcquire(ctx context.Context, key string, ttl time.Duration) (bool, error)
}
