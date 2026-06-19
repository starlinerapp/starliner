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
	organizationId int64,
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

	runner, err := qtx.CreateRunner(ctx, sql.NullInt64{Int64: organizationId, Valid: true})
	if err != nil {
		return nil, err
	}

	_, err = qtx.CreateRunnerRegistrationToken(ctx, sqlc.CreateRunnerRegistrationTokenParams{
		OrganizationID: sql.NullInt64{Int64: organizationId, Valid: true},
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

	return &entity.Runner{
		Id:                runner.ID,
		OrganizationId:    runner.OrganizationID.Int64,
		Name:              mapper.ToPtrFromNullString(runner.Name),
		Status:            runner.Status,
		Labels:            runner.Labels,
		MaxConcurrentJobs: runner.MaxConcurrentJobs,
		DisabledAt:        mapper.ToPtrFromNullTime(runner.DisabledAt),
		CreatedAt:         runner.CreatedAt,
		UpdatedAt:         runner.UpdatedAt,
	}, nil
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
