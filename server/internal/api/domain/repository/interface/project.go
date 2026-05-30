package interfaces

import (
	"context"

	"starliner.app/internal/api/domain/entity"
)

type ProjectRepository interface {
	CreateProjectWithEnvironment(
		ctx context.Context,
		projectName string,
		namespace string,
		environmentName string,
		environmentSlug string,
		organizationId int64,
		clusterId int64,
	) (*entity.Project, error)
	GetProject(ctx context.Context, projectId int64, userId int64) (*entity.Project, error)
	DeleteProject(ctx context.Context, projectId int64, userId int64) error
	GetProjectCluster(ctx context.Context, projectId int64, userId int64) (*entity.ProjectCluster, error)
	GetProjectEnvironments(ctx context.Context, projectId int64, userId int64) ([]*entity.Environment, error)
	GetProjectPreviewEnvironmentEnabled(ctx context.Context, projectId int64, userId int64) (bool, error)
	ToggleProjectPreviewEnvironmentEnabled(ctx context.Context, projectId int64, userId int64) (bool, error)
	GetProjectProductionEnvironmentsByRepositoryUrl(ctx context.Context, repositoryUrl string) ([]*entity.Environment, error)
	GetTeamProjectIds(ctx context.Context, teamId int64) ([]int64, error)
	GetProjectEnvironmentsByProjectId(ctx context.Context, projectId int64) ([]*entity.Environment, error)
	DeleteProjectsByTeamId(ctx context.Context, teamId int64) error
	GetProjectIdsByClusterId(ctx context.Context, clusterId int64) ([]int64, error)
	DeleteProjectsByClusterId(ctx context.Context, clusterId int64) error
}
