package repository

import (
	"context"
	"starliner.app/pkg/db/sqlc"
	"starliner.app/pkg/domain"
)

type EnvironmentRepository struct {
	queries *sqlc.Queries
}

func NewEnvironmentRepository(queries *sqlc.Queries) *EnvironmentRepository {
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
