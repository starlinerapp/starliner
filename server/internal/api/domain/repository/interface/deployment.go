package interfaces

import (
	"context"
	"starliner.app/internal/api/domain/entity"
)

type DeploymentRepository interface {
	CreateDatabaseDeployment(
		ctx context.Context,
		name string,
		port string,
		username string,
		password string,
		environmentId int64,
	) (deployment *entity.DatabaseDeployment, err error)

	GetUserDeployment(ctx context.Context, userId int64, deploymentId int64) (*entity.Deployment, error)
}
