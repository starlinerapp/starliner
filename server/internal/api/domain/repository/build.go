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

func (br *BuildRepository) CreateBuild(ctx context.Context, deploymentId int64) (*entity.Build, error) {
	b, err := br.queries.CreateBuild(ctx, utils.NullInt64FromPtr(&deploymentId))
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

func (br *BuildRepository) UpdateBuild(ctx context.Context, id int64, status value.BuildStatus, logs string) error {
	return br.queries.UpdateBuildInformation(ctx, sqlc.UpdateBuildInformationParams{
		Status: sqlc.BuildStatus(status),
		Logs:   utils.NullStringFromPtr(&logs),
		ID:     id,
	})
}

func (br *BuildRepository) GetBuildLogs(ctx context.Context, userId int64, buildId int64) (*string, error) {
	res, err := br.queries.GetBuildLogs(ctx, sqlc.GetBuildLogsParams{
		BuildID: buildId,
		UserID:  userId,
	})
	return utils.PtrFromNullString(res), err
}
