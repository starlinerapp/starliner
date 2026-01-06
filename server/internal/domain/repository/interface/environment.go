package interfaces

import (
	"context"
	"starliner.app/internal/domain/entity"
)

type EnvironmentRepository interface {
	CreateEnvironment(ctx context.Context, name string, slug string, projectId int64) (*entity.Environment, error)
}
