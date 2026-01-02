package repository

import (
	"context"
	"fmt"
	"starliner.app/internal/domain"
	"starliner.app/internal/infrastructure/db/sqlc"
	interfaces "starliner.app/internal/repository/interface"
)

type ProjectRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.ProjectRepository = (*ProjectRepository)(nil)

func NewProjectRepository(queries *sqlc.Queries) interfaces.ProjectRepository {
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

func (pr *ProjectRepository) GetProject(ctx context.Context, projectId int64, userId int64) (*domain.Project, error) {
	rows, err := pr.queries.GetProject(ctx, sqlc.GetProjectParams{
		ID:      projectId,
		OwnerID: userId,
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("project not found")
	}

	project := &domain.Project{
		Id:             rows[0].ID,
		Name:           rows[0].Name,
		OrganizationId: rows[0].OrganizationID,
		Environments:   make([]domain.Environment, 0, len(rows)),
	}

	for _, row := range rows {
		project.Environments = append(project.Environments, domain.Environment{
			Id:   row.EnvironmentID,
			Slug: row.EnvironmentSlug,
			Name: row.EnvironmentName,
		})
	}

	return project, nil
}
