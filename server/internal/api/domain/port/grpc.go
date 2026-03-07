package port

import (
	"context"
	"io"
)

type GrpcClient interface {
	StreamLogs(
		ctx context.Context,
		namespace string,
		releaseName string,
		kubeconfigBase64 string,
		w io.Writer,
	) error
}
