package port

import (
	"context"
	"io"
	"time"
)

type TerminalSize struct {
	Columns uint16
	Rows    uint16
}

type SSH interface {
	WaitForSSH(ip string, timeout time.Duration) error
	OpenTTY(
		ctx context.Context,
		user string,
		ip string,
		pemKey []byte,
		terminalSize <-chan TerminalSize,
	) (stdin io.WriteCloser, stdout io.ReadCloser, err error)
}
