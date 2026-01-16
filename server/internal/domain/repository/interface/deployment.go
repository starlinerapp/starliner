package interfaces

import (
	"context"
	"starliner.app/internal/domain/entity"
)

type DeploymentRepository interface {
	CreateDeployment(ctx context.Context, name string, environmentId int64) (deployment *entity.Deployment, err error)
}
