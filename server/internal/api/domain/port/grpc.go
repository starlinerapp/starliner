package port

import (
	"context"
	"io"
)

type TerminalSize struct {
	Columns int
	Rows    int
}

type BuilderClient interface {
	StreamBuildLogs(ctx context.Context, buildId int64, w io.Writer) error
}

type ClusterClient interface {
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

type ProvisionerClient interface {
	OpenTTY(
		ctx context.Context,
		ip string,
		user string,
		pemKey []byte,
		stdin io.Reader,
		stdout io.Writer,
		sizes <-chan TerminalSize,
	) error

	StreamProvisioningLogs(
		ctx context.Context,
		clusterId int64,
		w io.Writer,
	) error
}
