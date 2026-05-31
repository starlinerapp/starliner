package port

import (
	"context"
	"io"

	"starliner.app/internal/cluster/domain/value"
)

type Logs interface {
	StreamLogs(
		ctx context.Context,
		source value.LogSource,
		environmentNamespace string,
		releaseName string,
		kubeconfigBase64 string,
	) (io.ReadCloser, error)
}

type LogPublisher interface {
	PublishLogChunk(ctx context.Context, namespace string, releaseName string, data []byte) error
}
