package interfaces

import (
	"context"
	"starliner.app/internal/domain/entity"
)

type EnvironmentRepository interface {
	CreateEnvironment(ctx context.Context, name string, slug string, projectId int64) (*entity.Environment, error)
	GetEnvironmentAuthorizedUsers(ctx context.Context, clusterId int64) (users []int64, err error)
	GetEnvironmentCluster(ctx context.Context, environmentId int64) (*entity.Cluster, error)
}
