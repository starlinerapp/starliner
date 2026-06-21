package repository

import (
	"context"
	"database/sql"
	"time"

	"starliner.app/internal/api/domain/entity"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/infrastructure/postgres/mapper"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
)

type RunnerRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

var _ interfaces.RunnerRepository = (*RunnerRepository)(nil)

func NewRunnerRepository(db *sql.DB, queries *sqlc.Queries) interfaces.RunnerRepository {
	return &RunnerRepository{
		db:      db,
		queries: queries,
	}
}

func (rr *RunnerRepository) CreateRunnerWithRegistrationToken(
	ctx context.Context,
	organizationId *int64,
	tokenHash string,
	expiresAt time.Time,
) (*entity.Runner, error) {
	tx, err := rr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := rr.queries.WithTx(tx)

	orgID := mapper.ToNullInt64FromPtr(organizationId)

	runner, err := qtx.CreateRunner(ctx, orgID)
	if err != nil {
		return nil, err
	}

	_, err = qtx.CreateRunnerRegistrationToken(ctx, sqlc.CreateRunnerRegistrationTokenParams{
		OrganizationID: orgID,
		TokenHash:      tokenHash,
		ExpiresAt:      expiresAt,
		RunnerID:       sql.NullInt64{Int64: runner.ID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return mapRunnerRow(runner), nil
}

func (rr *RunnerRepository) RegisterRunner(
	ctx context.Context,
	tokenHash string,
	name string,
	labels []string,
	maxConcurrentJobs int32,
) error {
	tx, err := rr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := rr.queries.WithTx(tx)

	regToken, err := qtx.UseRunnerRegistrationToken(ctx, tokenHash)
	if err != nil {
		return err
	}

	if !regToken.RunnerID.Valid {
		return sql.ErrNoRows
	}

	_, err = qtx.UpdateRunnerOnRegistration(ctx, sqlc.UpdateRunnerOnRegistrationParams{
		ID:                regToken.RunnerID.Int64,
		Name:              sql.NullString{String: name, Valid: true},
		Labels:            labels,
		MaxConcurrentJobs: maxConcurrentJobs,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (rr *RunnerRepository) GetOrganizationRunners(
	ctx context.Context,
	organizationId int64,
) ([]*entity.Runner, error) {
	rows, err := rr.queries.GetOrganizationRunners(ctx, sql.NullInt64{Int64: organizationId, Valid: true})
	if err != nil {
		return nil, err
	}

	runners := make([]*entity.Runner, len(rows))
	for i, row := range rows {
		runners[i] = mapRunnerRow(row)
	}

	return runners, nil
}

func (rr *RunnerRepository) GetGlobalRunners(ctx context.Context) ([]*entity.Runner, error) {
	rows, err := rr.queries.GetGlobalRunners(ctx)
	if err != nil {
		return nil, err
	}

	runners := make([]*entity.Runner, len(rows))
	for i, row := range rows {
		runners[i] = mapRunnerRow(row)
	}

	return runners, nil
}

func (rr *RunnerRepository) GetRunnerById(
	ctx context.Context,
	runnerId int64,
) (*entity.Runner, error) {
	row, err := rr.queries.GetRunnerById(ctx, runnerId)
	if err != nil {
		return nil, err
	}

	return mapRunnerRow(row), nil
}

func mapRunnerRow(row sqlc.Runner) *entity.Runner {
	return &entity.Runner{
		Id:                row.ID,
		OrganizationId:    mapper.ToPtrFromNullInt64(row.OrganizationID),
		Name:              mapper.ToPtrFromNullString(row.Name),
		Status:            row.Status,
		Labels:            row.Labels,
		MaxConcurrentJobs: row.MaxConcurrentJobs,
		DisabledAt:        mapper.ToPtrFromNullTime(row.DisabledAt),
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}

func (rr *RunnerRepository) GetRunnerIdByRegistrationToken(
	ctx context.Context,
	tokenHash string,
) (int64, error) {
	runnerID, err := rr.queries.GetRunnerIdByRegistrationToken(ctx, tokenHash)
	if err != nil {
		return 0, err
	}

	if !runnerID.Valid {
		return 0, sql.ErrNoRows
	}

	return runnerID.Int64, nil
}

func (rr *RunnerRepository) ResolveRunnerByRegistrationToken(
	ctx context.Context,
	tokenHash string,
) (int64, *int64, error) {
	row, err := rr.queries.ResolveRunnerByRegistrationToken(ctx, tokenHash)
	if err != nil {
		return 0, nil, err
	}

	if !row.RunnerID.Valid {
		return 0, nil, sql.ErrNoRows
	}

	return row.RunnerID.Int64, mapper.ToPtrFromNullInt64(row.OrganizationID), nil
}

func (rr *RunnerRepository) UpdateRunnerStatus(
	ctx context.Context,
	runnerId int64,
	status string,
) error {
	return rr.queries.UpdateRunnerStatus(ctx, sqlc.UpdateRunnerStatusParams{
		ID:     runnerId,
		Status: status,
	})
}

func (rr *RunnerRepository) GetRunnerByOrganization(
	ctx context.Context,
	runnerId int64,
	organizationId int64,
) (*entity.Runner, error) {
	row, err := rr.queries.GetRunnerByOrganization(ctx, sqlc.GetRunnerByOrganizationParams{
		ID:             runnerId,
		OrganizationID: sql.NullInt64{Int64: organizationId, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return mapRunnerRow(row), nil
}

func (rr *RunnerRepository) DeleteRunner(
	ctx context.Context,
	runnerId int64,
	organizationId int64,
) error {
	return rr.queries.DeleteRunner(ctx, sqlc.DeleteRunnerParams{
		ID:             runnerId,
		OrganizationID: sql.NullInt64{Int64: organizationId, Valid: true},
	})
}

func (rr *RunnerRepository) DeleteGlobalRunner(ctx context.Context, runnerId int64) error {
	return rr.queries.DeleteGlobalRunner(ctx, runnerId)
}
