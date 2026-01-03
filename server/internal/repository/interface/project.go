package interfaces

import (
	"context"
	"starliner.app/internal/domain"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, name string, organizationId int64, clusterId int64) (*domain.Project, error)
	GetProject(ctx context.Context, projectId int64, userId int64) (*domain.Project, error)
}
