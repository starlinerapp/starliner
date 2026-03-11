package port

import (
	"context"
	"io"
)

type TerminalSize struct {
	Columns uint16
	Rows    uint16
}

type TTY interface {
	Open(
		ctx context.Context,
		namespace string,
		releaseName string,
		kubeconfigBase64 string,
		terminalSize <-chan TerminalSize,
	) (stdin io.WriteCloser, stdout io.ReadCloser, err error)
}
