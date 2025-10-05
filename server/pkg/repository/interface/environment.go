package interfaces

import (
	"context"
	"starliner.app/pkg/domain"
)

type EnvironmentRepository interface {
	CreateEnvironment(ctx context.Context, name string, projectId int64) (*domain.Environment, error)
}
