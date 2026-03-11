package port

import (
	"context"
	"io"
)

type TerminalSize struct {
	Columns int
	Rows    int
}

type GrpcClient interface {
	StreamLogs(
		ctx context.Context,
		namespace string,
		releaseName string,
		kubeconfigBase64 string,
		w io.Writer,
	) error

	OpenTTY(
		ctx context.Context,
		namespace string,
		releaseName string,
		kubeconfigBase64 string,
		stdin io.Reader,
		stdout io.Writer,
		sizes <-chan TerminalSize,
	) error
}
