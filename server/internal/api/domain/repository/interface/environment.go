package interfaces

import (
	"context"
	"starliner.app/internal/api/domain/entity"
)

type EnvironmentRepository interface {
	CreateEnvironment(ctx context.Context, name string, slug string, projectId int64) (*entity.Environment, error)
	GetEnvironmentAuthorizedUsers(ctx context.Context, clusterId int64) (users []int64, err error)
	GetEnvironmentCluster(ctx context.Context, environmentId int64) (*entity.Cluster, error)
	GetEnvironmentIngressDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.IngressDeployment, error)
	GetEnvironmentImageDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.ImageDeployment, error)
	GetEnvironmentDatabaseDeployments(ctx context.Context, environmentId int64, userId int64) (deployments []*entity.DatabaseDeployment, err error)
	GetEnvironmentDeploymentByName(ctx context.Context, name string, environmentId int64) (*entity.Deployment, error)
}
