package interfaces

import (
	"context"
	"starliner.app/internal/domain"
)

type EnvironmentRepository interface {
	CreateEnvironment(ctx context.Context, name string, slug string, projectId int64) (*domain.Environment, error)
}
