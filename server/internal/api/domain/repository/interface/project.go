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
}
