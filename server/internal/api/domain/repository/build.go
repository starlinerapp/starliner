package repository

import (
	"context"
	"database/sql"

	"starliner.app/internal/api/domain/entity"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
	"starliner.app/internal/api/infrastructure/postgres/utils"
)

type BuildRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

var _ interfaces.BuildRepository = (*BuildRepository)(nil)

func NewBuildRepository(db *sql.DB, queries *sqlc.Queries) interfaces.BuildRepository {
	return &BuildRepository{db: db, queries: queries}
}

func (br *BuildRepository) CreateBuild(ctx context.Context, deploymentId int64, source string, args []*value.Arg) (*entity.Build, error) {
	tx, err := br.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := br.queries.WithTx(tx)

	b, err := qtx.CreateBuild(ctx, sqlc.CreateBuildParams{
		DeploymentID: utils.NullInt64FromPtr(&deploymentId),
		Source:       source,
	})
	if err != nil {
		return nil, err
	}

	resultArgs := make([]*entity.Arg, len(args))
	for i, a := range args {
		arg, err := qtx.CreateBuildArg(ctx, sqlc.CreateBuildArgParams{
			BuildID: b.ID,
			Name:    a.Name,
			Value:   a.Value,
		})
		if err != nil {
			return nil, err
		}

		resultArgs[i] = &entity.Arg{
			Name:  arg.Name,
			Value: arg.Value,
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &entity.Build{
		Id:           b.ID,
		DeploymentId: utils.PtrFromNullInt64(b.DeploymentID),
		Status:       entity.BuildStatus(b.Status),
		Logs:         utils.PtrFromNullString(b.Logs),
		Args:         resultArgs,
	}, nil
}

func (br *BuildRepository) GetBuildArgs(ctx context.Context, buildId int64) ([]*entity.Arg, error) {
	args, err := br.queries.GetBuildArgs(ctx, buildId)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Arg, len(args))
	for i, a := range args {
		result[i] = &entity.Arg{
			Name:  a.Name,
			Value: a.Value,
		}
	}
	return result, nil
}

func (br *BuildRepository) GetLatestBuildArgs(ctx context.Context, deploymentId int64) ([]*entity.Arg, error) {
	buildId, err := br.queries.GetLatestBuildIdByDeploymentId(ctx, utils.NullInt64FromPtr(&deploymentId))
	if err != nil {
		return nil, err
	}

	return br.GetBuildArgs(ctx, buildId)
}

func (br *BuildRepository) UpdateBuild(ctx context.Context, id int64, status value.BuildStatus, commitHash *string, logs string) error {
	return br.queries.UpdateBuildInformation(ctx, sqlc.UpdateBuildInformationParams{
		Status:     sqlc.BuildStatus(status),
		CommitHash: utils.NullStringFromPtr(commitHash),
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
