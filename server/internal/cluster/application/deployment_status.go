package application

import (
	"context"
	"io"

	"starliner.app/internal/cluster/domain/port"
)

type DeploymentStatusApplication struct {
	deploymentStatus port.DeploymentStatus
}

func NewDeploymentStatusApplication(deploymentStatus port.DeploymentStatus) *DeploymentStatusApplication {
	return &DeploymentStatusApplication{deploymentStatus: deploymentStatus}
}

func (a *DeploymentStatusApplication) StreamDeploymentStatusLogs(
	ctx context.Context,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
	commitHash string,
) (io.ReadCloser, error) {
	return a.deploymentStatus.StreamDeploymentStatusLogs(
		ctx,
		namespace,
		releaseName,
		kubeconfigBase64,
		commitHash,
	)
}
