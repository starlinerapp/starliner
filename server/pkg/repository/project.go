package repository

import (
	"context"
	"starliner.app/pkg/db/sqlc"
	"starliner.app/pkg/domain"
)

type ProjectRepository struct {
	queries *sqlc.Queries
}

func NewProjectRepository(queries *sqlc.Queries) *ProjectRepository {
	return &ProjectRepository{queries: queries}
}

func (pr *ProjectRepository) CreateProject(ctx context.Context, name string, organizationId int64) (*domain.Project, error) {
	project, err := pr.queries.CreateProject(ctx, sqlc.CreateProjectParams{
		Name:           name,
		OrganizationID: organizationId,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Project{
		Id:             project.ID,
		Name:           project.Name,
		OrganizationId: project.OrganizationID,
	}, nil
}
