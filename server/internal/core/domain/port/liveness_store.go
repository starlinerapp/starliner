package port

import (
	"context"
	"time"
)

type LivenessStore interface {
	MarkAlive(ctx context.Context, key string, ttl time.Duration) error
}
