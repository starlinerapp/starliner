package repository

import (
	"context"
	"database/sql"
	"fmt"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
	"starliner.app/internal/api/infrastructure/postgres/utils"
)

type ProjectRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

var _ interfaces.ProjectRepository = (*ProjectRepository)(nil)

func NewProjectRepository(db *sql.DB, queries *sqlc.Queries) interfaces.ProjectRepository {
	return &ProjectRepository{db: db, queries: queries}
}

func (pr *ProjectRepository) CreateProjectWithEnvironment(
	ctx context.Context,
	projectName string,
	namespace string,
	environmentName string,
	environmentSlug string,
	organizationId int64,
	clusterId int64,
) (*entity.Project, error) {
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := pr.queries.WithTx(tx)
	project, err := qtx.CreateProject(ctx, sqlc.CreateProjectParams{
		Name:           projectName,
		OrganizationID: organizationId,
		ClusterID:      utils.NullInt64FromPtr(&clusterId),
	})
	if err != nil {
		return nil, err
	}

	env, err := qtx.CreateEnvironment(ctx, sqlc.CreateEnvironmentParams{
		Name:      environmentName,
		Namespace: namespace,
		Slug:      environmentSlug,
		ProjectID: project.ID,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &entity.Project{
		Id:             project.ID,
		Name:           project.Name,
		OrganizationId: project.OrganizationID,
		ClusterId:      utils.PtrFromNullInt64(project.ClusterID),
		Environments: []*entity.Environment{
			{
				Id:   env.ID,
				Slug: env.Slug,
				Name: env.Name,
			},
		},
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

func (pr *ProjectRepository) GetProjectCluster(ctx context.Context, projectId int64, userId int64) (*entity.ProjectCluster, error) {
	row, err := pr.queries.GetProjectCluster(ctx, sqlc.GetProjectClusterParams{
		ID:      projectId,
		OwnerID: userId,
	})
	if err != nil {
		return nil, err
	}
	
	return &entity.ProjectCluster{
		ClusterId:   row.ID,
		ClusterName: row.Name,
	}, nil
}
