package repository

import (
	"context"
	"fmt"
	"starliner.app/internal/domain/entity"
	"starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/infrastructure/postgres/sqlc"
	"starliner.app/internal/infrastructure/postgres/utils"
)

type ProjectRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.ProjectRepository = (*ProjectRepository)(nil)

func NewProjectRepository(queries *sqlc.Queries) interfaces.ProjectRepository {
	return &ProjectRepository{queries: queries}
}

func (pr *ProjectRepository) CreateProject(ctx context.Context, name string, organizationId int64, clusterId int64) (*entity.Project, error) {
	project, err := pr.queries.CreateProject(ctx, sqlc.CreateProjectParams{
		Name:           name,
		OrganizationID: organizationId,
		ClusterID:      utils.NullInt64FromPtr(&clusterId),
	})
	if err != nil {
		return nil, err
	}

	return &entity.Project{
		Id:             project.ID,
		Name:           project.Name,
		OrganizationId: project.OrganizationID,
		ClusterId:      utils.PtrFromNullInt64(project.ClusterID),
	}, nil
}

func (pr *ProjectRepository) GetProject(ctx context.Context, projectId int64, userId int64) (*entity.Project, error) {
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

	project := &entity.Project{
		Id:             rows[0].ID,
		Name:           rows[0].Name,
		OrganizationId: rows[0].OrganizationID,
		ClusterId:      utils.PtrFromNullInt64(rows[0].ClusterID),
		Environments:   make([]*entity.Environment, 0, len(rows)),
	}

	for _, row := range rows {
		project.Environments = append(project.Environments, &entity.Environment{
			Id:   row.EnvironmentID,
			Slug: row.EnvironmentSlug,
			Name: row.EnvironmentName,
		})
	}

	return project, nil
}

func (pr *ProjectRepository) DeleteProject(ctx context.Context, projectId int64, userId int64) error {
	return pr.queries.DeleteProject(ctx, sqlc.DeleteProjectParams{
		ID:      projectId,
		OwnerID: userId,
	})
}
