package application

import (
	"context"
	"io"
	"starliner.app/internal/provisioner/domain/port"
)

type TTYApplication struct {
	ssh port.SSH
}

func NewTTYApplication(ssh port.SSH) *TTYApplication {
	return &TTYApplication{ssh: ssh}
}

func (ta *TTYApplication) OpenTTY(
	ctx context.Context,
	user string,
	ip string,
	pemKey []byte,
	terminalSize <-chan port.TerminalSize,
) (stdin io.WriteCloser, stdout io.ReadCloser, err error) {
	return ta.ssh.OpenTTY(ctx, user, ip, pemKey, terminalSize)
}
