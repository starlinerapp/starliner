package repository

import (
	"context"
	"database/sql"
	"errors"
	"starliner.app/internal/api/domain/entity"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
)

type GithubAppRepository struct {
	queries *sqlc.Queries
}

func NewGithubAppRepository(queries *sqlc.Queries) interfaces.GithubAppRepository {
	return &GithubAppRepository{
		queries: queries,
	}
}

func (gr *GithubAppRepository) CreateGithubApp(ctx context.Context, installationID int64, organizationId int64) (*entity.GithubApp, error) {
	ghApp, err := gr.queries.CreateGithubApp(ctx, sqlc.CreateGithubAppParams{
		InstallationID: installationID,
		OrganizationID: organizationId,
	})
	if err != nil {
		return nil, err
	}

	return &entity.GithubApp{
		ID:             ghApp.ID,
		InstallationID: ghApp.InstallationID,
		OrganizationID: ghApp.OrganizationID,
		CreatedAt:      ghApp.CreatedAt,
	}, nil
}

func (gr *GithubAppRepository) GetOrganizationGithubApp(ctx context.Context, organizationId int64) (*entity.GithubApp, error) {
	ghApp, err := gr.queries.GetOrganizationGithubApp(ctx, organizationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity.GithubApp{
		ID:             ghApp.ID,
		InstallationID: ghApp.InstallationID,
		OrganizationID: ghApp.OrganizationID,
		CreatedAt:      ghApp.CreatedAt,
	}, nil
}

func (gr *GithubAppRepository) GetEnvironmentGithubApp(ctx context.Context, environmentId int64) (*entity.GithubApp, error) {
	ghApp, err := gr.queries.GetEnvironmentGithubApp(ctx, environmentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity.GithubApp{
		ID:             ghApp.ID,
		InstallationID: ghApp.InstallationID,
		OrganizationID: ghApp.OrganizationID,
		CreatedAt:      ghApp.CreatedAt,
	}, nil
}

func (gr *GithubAppRepository) DeleteGithubApp(ctx context.Context, installationId int64) error {
	return gr.queries.DeleteGithubAppByInstallationId(ctx, installationId)
}
