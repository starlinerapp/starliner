package port

import (
	"context"
	"io"
)

type DeploymentStatus interface {
	StreamDeploymentStatusLogs(
		ctx context.Context,
		namespace string,
		releaseName string,
		kubeconfigBase64 string,
		commitHash string,
	) (io.ReadCloser, error)
}
