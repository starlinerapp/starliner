package interfaces

import (
	"context"

	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/value"
)

type BuildRepository interface {
	CreateBuild(ctx context.Context, deploymentId int64, source string, args []*value.Arg) (*entity.Build, error)
	UpdateBuild(ctx context.Context, id int64, status value.BuildStatus, commitHash *string, logs string) error
	GetBuildLogs(ctx context.Context, userId int64, buildId int64) (*string, error)
	GetBuildArgs(ctx context.Context, buildId int64) ([]*entity.Arg, error)
	GetLatestBuildArgs(ctx context.Context, deploymentId int64) ([]*entity.Arg, error)
}
