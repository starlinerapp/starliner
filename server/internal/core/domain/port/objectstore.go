package port

import (
	"context"
	"io"
)

type ObjectStore interface {
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	CreateBuckets(ctx context.Context) error
}
