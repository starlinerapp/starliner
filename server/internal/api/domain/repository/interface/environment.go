package interfaces

import (
	"context"

	"starliner.app/internal/api/domain/entity"
)

type EnvironmentRepository interface {
	CreateEnvironment(ctx context.Context, name string, namespace string, slug string, projectId int64) (*entity.Environment, error)
	GetEnvironmentById(ctx context.Context, environmentId int64) (*entity.Environment, error)
	DuplicateEnvironment(ctx context.Context, userId int64, name string, namespace string, slug string, projectId int64, sourceEnvironmentId int64, uniqueIdentifier string) (*entity.Environment, error)
	GetEnvironmentAuthorizedUsers(ctx context.Context, clusterId int64) (users []int64, err error)
	GetEnvironmentCluster(ctx context.Context, environmentId int64) (*entity.Cluster, error)
	GetUserEnvironmentGitDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.GitDeployment, error)
	GetEnvironmentIngressDeploymentByName(ctx context.Context, environmentId int64, name string) (*entity.IngressDeployment, error)
	GetUserEnvironmentIngressDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.IngressDeployment, error)
	GetUserEnvironmentImageDeployments(ctx context.Context, environmentId int64, userId int64) ([]*entity.ImageDeployment, error)
	GetUserEnvironmentDatabaseDeployments(ctx context.Context, environmentId int64, userId int64) (deployments []*entity.DatabaseDeployment, err error)
	GetEnvironmentDeploymentByName(ctx context.Context, name string, environmentId int64) (*entity.Deployment, error)
	GetEnvironmentGitDeploymentBuilds(ctx context.Context, environmentId int64) ([]*entity.GitDeploymentBuild, error)
	GetEnvironmentBranch(ctx context.Context, environmentId int64) (string, error)
	UpdateEnvironmentBranch(ctx context.Context, environmentId int64, branch string) error
}
