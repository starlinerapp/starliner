package interfaces

import (
	"context"

	"starliner.app/internal/api/domain/entity"
)

type EnvironmentRepository interface {
	CreateEnvironment(ctx context.Context, name string, namespace string, slug string, projectId int64) (*entity.Environment, error)
	DeleteEnvironment(ctx context.Context, environmentId int64) error
	GetEnvironmentById(ctx context.Context, environmentId int64) (*entity.Environment, error)
	GetPreviewEnvironment(ctx context.Context, gitHubRepositoryId int64, prNumber int) (*entity.PreviewEnvironment, error)
	GetEnvironmentProject(ctx context.Context, environmentId int64) (*entity.Project, error)
	CloneEnvironment(ctx context.Context, name string, namespace string, slug string, projectId int64, sourceEnvironmentId int64, uniqueIdentifier string, connectedBranch *string, preview *entity.PreviewCloneMetadata) (*entity.Environment, error)
	GetEnvironmentAuthorizedUsers(ctx context.Context, clusterId int64) (users []int64, err error)
	GetEnvironmentCluster(ctx context.Context, environmentId int64) (*entity.Cluster, error)
	GetEnvironmentOrganization(ctx context.Context, environmentId int64) (*entity.Organization, error)
	GetEnvironmentIngressDeploymentByName(ctx context.Context, environmentId int64, name string) (*entity.IngressDeployment, error)
	GetEnvironmentGitDeployments(ctx context.Context, environmentId int64) ([]*entity.GitDeployment, error)
	GetEnvironmentIngressDeployments(ctx context.Context, environmentId int64) ([]*entity.IngressDeployment, error)
	GetEnvironmentImageDeployments(ctx context.Context, environmentId int64) ([]*entity.ImageDeployment, error)
	GetEnvironmentDatabaseDeployments(ctx context.Context, environmentId int64) (deployments []*entity.DatabaseDeployment, err error)
	GetEnvironmentDeploymentByName(ctx context.Context, name string, environmentId int64) (*entity.Deployment, error)
	GetEnvironmentGitDeploymentBuilds(ctx context.Context, environmentId int64) ([]*entity.GitDeploymentBuild, error)
	GetEnvironmentIngressDeploymentBuilds(ctx context.Context, environmentId int64) ([]*entity.GitDeploymentBuild, error)
	GetEnvironmentImageDeploymentBuilds(ctx context.Context, environmentId int64) ([]*entity.GitDeploymentBuild, error)
	GetEnvironmentDatabaseDeploymentBuilds(ctx context.Context, environmentId int64) ([]*entity.GitDeploymentBuild, error)
	GetEnvironmentBranch(ctx context.Context, environmentId int64) (string, error)
	UpdateEnvironmentBranch(ctx context.Context, environmentId int64, branch string) error
}
