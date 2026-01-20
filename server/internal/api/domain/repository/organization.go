package repository

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	sqlc2 "starliner.app/internal/api/infrastructure/postgres/sqlc"
)

type OrganizationRepository struct {
	queries *sqlc2.Queries
}

var _ interfaces.OrganizationRepository = (*OrganizationRepository)(nil)

func NewOrganizationRepository(queries *sqlc2.Queries) interfaces.OrganizationRepository {
	return &OrganizationRepository{queries: queries}
}

func (or *OrganizationRepository) CreateOrganization(ctx context.Context, name string, slug string, ownerID int64) (*entity.Organization, error) {
	organization, err := or.queries.CreateOrganization(ctx, sqlc2.CreateOrganizationParams{
		Name:    name,
		Slug:    slug,
		OwnerID: ownerID,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Organization{
		Id:      organization.ID,
		Name:    organization.Name,
		Slug:    organization.Slug,
		OwnerId: organization.OwnerID,
	}, nil
}

func (or *OrganizationRepository) GetOrganization(ctx context.Context, id int64) (*entity.Organization, error) {
	organization, err := or.queries.GetOrganization(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.Organization{
		Id:      organization.ID,
		Name:    organization.Name,
		Slug:    organization.Slug,
		OwnerId: organization.OwnerID,
	}, nil
}

func (or *OrganizationRepository) GetUserOrganizations(ctx context.Context, userId int64) ([]*entity.Organization, error) {
	organizations, err := or.queries.GetUserOrganizations(ctx, userId)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Organization, len(organizations))
	for i, organization := range organizations {
		result[i] = &entity.Organization{
			Id:      organization.ID,
			Name:    organization.Name,
			Slug:    organization.Slug,
			OwnerId: organization.OwnerID,
		}
	}

	return result, nil
}

func (or *OrganizationRepository) GetOrganizationProjects(ctx context.Context, organizationID int64) ([]*entity.Project, error) {
	rows, err := or.queries.GetOrganizationProjects(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	projectsMap := make(map[int64]*entity.Project)
	for _, row := range rows {
		proj, exists := projectsMap[row.ID]
		if !exists {
			proj = &entity.Project{
				Id:             row.ID,
				Name:           row.Name,
				Environments:   []*entity.Environment{},
				OrganizationId: row.OrganizationID,
			}
			projectsMap[proj.Id] = proj
		}

		proj.Environments = append(proj.Environments, &entity.Environment{
			Id:   row.EnvironmentID,
			Slug: row.EnvironmentSlug,
			Name: row.EnvironmentName,
		})
	}

	projects := make([]*entity.Project, 0, len(projectsMap))
	for _, p := range projectsMap {
		projects = append(projects, p)
	}

	return projects, nil
}

func (or *OrganizationRepository) GetOrganizationClusters(ctx context.Context, organizationID int64) ([]*entity.Cluster, error) {
	rows, err := or.queries.GetOrganizationClusters(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if rows == nil {
		return []*entity.Cluster{}, nil
	}

	clusters := make([]*entity.Cluster, 0, len(rows))
	for _, c := range rows {
		clusters = append(clusters, &entity.Cluster{
			Id:             c.ID,
			Name:           c.Name,
			OrganizationId: c.OrganizationID,
		})
	}

	return clusters, nil
}
