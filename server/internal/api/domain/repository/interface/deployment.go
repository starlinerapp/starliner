package _interface

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
	) (deployment *entity.Deployment, err error)
}
