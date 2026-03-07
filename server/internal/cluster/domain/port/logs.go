package port

import (
	"context"
	"io"
)

type Logs interface {
	StreamLogs(ctx context.Context, namespace string, releaseName string, kubeconfigBase64 string) (io.ReadCloser, error)
}
