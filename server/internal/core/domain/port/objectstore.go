package port

import (
	"context"
	"io"
)

type ObjectStore interface {
	CreateBuckets(ctx context.Context) error
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	UploadDirAsTarGz(ctx context.Context, pathToDir string, key string) (string, error)
}
