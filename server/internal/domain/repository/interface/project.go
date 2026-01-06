package interfaces

import (
	"context"
	"starliner.app/internal/domain/entity"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, name string, organizationId int64, clusterId int64) (*entity.Project, error)
	GetProject(ctx context.Context, projectId int64, userId int64) (*entity.Project, error)
	DeleteProject(ctx context.Context, projectId int64, userId int64) error
}
