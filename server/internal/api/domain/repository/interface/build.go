package interfaces

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/value"
)

type BuildRepository interface {
	CreateBuild(ctx context.Context, deploymentId int64, source string) (*entity.Build, error)
	UpdateBuild(ctx context.Context, id int64, status value.BuildStatus, commitHash *string, imageName *string, logs string) error
	GetBuildLogs(ctx context.Context, userId int64, buildId int64) (*string, error)
	GetLatestGitDeploymentBuild(ctx context.Context, environmentId int64, serviceName string) (*entity.GitDeploymentBuild, error)
}
