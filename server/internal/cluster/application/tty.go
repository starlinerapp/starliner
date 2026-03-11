package application

import (
	"context"
	"io"
	"starliner.app/internal/cluster/domain/port"
)

type TTYApplication struct {
	tty port.TTY
}

func NewTTYApplication(tty port.TTY) *TTYApplication {
	return &TTYApplication{tty: tty}
}

func (ta *TTYApplication) OpenTTY(
	ctx context.Context,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
	terminalSize <-chan port.TerminalSize,
) (stdin io.WriteCloser, stdout io.ReadCloser, err error) {
	return ta.tty.Open(ctx, namespace, releaseName, kubeconfigBase64, terminalSize)
}
