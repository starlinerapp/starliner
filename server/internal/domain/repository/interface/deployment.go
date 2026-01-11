package interfaces

import "context"

type DeploymentRepository interface {
	CreateDeployment(ctx context.Context, name string, environmentId int64) error
}
