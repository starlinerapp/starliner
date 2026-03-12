package interfaces

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/value"
)

type BuildRepository interface {
	CreateBuild(ctx context.Context, deploymentId int64) (*entity.Build, error)
	UpdateBuild(ctx context.Context, id int64, status value.BuildStatus, logs string) error
}
