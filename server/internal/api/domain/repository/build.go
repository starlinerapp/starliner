package repository

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
	"starliner.app/internal/api/infrastructure/postgres/utils"
)

type BuildRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.BuildRepository = (*BuildRepository)(nil)

func NewBuildRepository(queries *sqlc.Queries) interfaces.BuildRepository {
	return &BuildRepository{queries: queries}
}

func (br *BuildRepository) CreateBuild(ctx context.Context, deploymentId int64, source string) (*entity.Build, error) {
	b, err := br.queries.CreateBuild(ctx, sqlc.CreateBuildParams{
		DeploymentID: utils.NullInt64FromPtr(&deploymentId),
		Source:       source,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Build{
		Id:           b.ID,
		DeploymentId: utils.PtrFromNullInt64(b.DeploymentID),
		Status:       entity.BuildStatus(b.Status),
		Logs:         utils.PtrFromNullString(b.Logs),
	}, nil
}

func (br *BuildRepository) UpdateBuild(ctx context.Context, id int64, status value.BuildStatus, commitHash *string, imageName *string, logs string) error {
	return br.queries.UpdateBuildInformation(ctx, sqlc.UpdateBuildInformationParams{
		Status:     sqlc.BuildStatus(status),
		CommitHash: utils.NullStringFromPtr(commitHash),
		ImageName:  utils.NullStringFromPtr(imageName),
		Logs:       utils.NullStringFromPtr(&logs),
		ID:         id,
	})
}

func (br *BuildRepository) GetBuildLogs(ctx context.Context, userId int64, buildId int64) (*string, error) {
	res, err := br.queries.GetBuildLogs(ctx, sqlc.GetBuildLogsParams{
		BuildID: buildId,
		UserID:  userId,
	})
	return utils.PtrFromNullString(res), err
}

func (br *BuildRepository) GetLatestGitDeploymentBuild(ctx context.Context, environmentId int64, serviceName string) (*entity.GitDeploymentBuild, error) {
	build, err := br.queries.GetLatestGitDeploymentBuild(ctx, sqlc.GetLatestGitDeploymentBuildParams{
		EnvironmentID: utils.NullInt64FromPtr(&environmentId),
		Name:          serviceName,
	})
	if err != nil {
		return nil, err
	}

	return &entity.GitDeploymentBuild{
		BuildId:        build.BuildID,
		DeploymentId:   build.DeploymentID,
		DeploymentName: build.DeploymentName,
		ImageName:      utils.PtrFromNullString(build.ImageName),
		CommitHash:     utils.PtrFromNullString(build.CommitHash),
		Source:         build.Source,
		Status:         entity.BuildStatus(build.Status),
		GitUrl:         build.Url,
		ProjectPath:    build.ProjectPath,
		DockerfilePath: build.DockerfilePath,
		CreatedAt:      build.CreatedAt,
	}, nil
}
