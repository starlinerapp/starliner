package repository

import (
	"context"
	"starliner.app/internal/domain"
	"starliner.app/internal/infrastructure/db/sqlc"
	interfaces "starliner.app/internal/repository/interface"
)

type EnvironmentRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.EnvironmentRepository = (*EnvironmentRepository)(nil)

func NewEnvironmentRepository(queries *sqlc.Queries) interfaces.EnvironmentRepository {
	return &EnvironmentRepository{queries: queries}
}

func (er *EnvironmentRepository) CreateEnvironment(ctx context.Context, name string, slug string, projectId int64) (*domain.Environment, error) {
	env, err := er.queries.CreateEnvironment(ctx, sqlc.CreateEnvironmentParams{
		Name:      name,
		Slug:      slug,
		ProjectID: projectId,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Environment{
		Id:   env.ID,
		Slug: env.Slug,
		Name: env.Name,
	}, nil
}
