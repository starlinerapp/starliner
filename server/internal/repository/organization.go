package repository

import (
	"context"
	"starliner.app/internal/domain"
	"starliner.app/internal/infrastructure/db/sqlc"
	interfaces "starliner.app/internal/repository/interface"
)

type OrganizationRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.OrganizationRepository = (*OrganizationRepository)(nil)

func NewOrganizationRepository(queries *sqlc.Queries) interfaces.OrganizationRepository {
	return &OrganizationRepository{queries: queries}
}

func (or *OrganizationRepository) CreateOrganization(ctx context.Context, name string, slug string, ownerID int64) (*domain.Organization, error) {
	organization, err := or.queries.CreateOrganization(ctx, sqlc.CreateOrganizationParams{
		Name:    name,
		Slug:    slug,
		OwnerID: ownerID,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Organization{
		Id:      organization.ID,
		Name:    organization.Name,
		Slug:    organization.Slug,
		OwnerId: organization.OwnerID,
	}, nil
}

func (or *OrganizationRepository) GetOrganization(ctx context.Context, id int64) (*domain.Organization, error) {
	organization, err := or.queries.GetOrganization(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Organization{
		Id:      organization.ID,
		Name:    organization.Name,
		Slug:    organization.Slug,
		OwnerId: organization.OwnerID,
	}, nil
}

func (or *OrganizationRepository) GetUserOrganizations(ctx context.Context, userId int64) ([]*domain.Organization, error) {
	organizations, err := or.queries.GetUserOrganizations(ctx, userId)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Organization, len(organizations))
	for i, organization := range organizations {
		result[i] = &domain.Organization{
			Id:      organization.ID,
			Name:    organization.Name,
			Slug:    organization.Slug,
			OwnerId: organization.OwnerID,
		}
	}

	return result, nil
}

func (or *OrganizationRepository) GetOrganizationProjects(ctx context.Context, organizationID int64) ([]*domain.Project, error) {
	rows, err := or.queries.GetOrganizationProjects(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	projectsMap := make(map[int64]*domain.Project)
	for _, row := range rows {
		proj, exists := projectsMap[row.ID]
		if !exists {
			proj = &domain.Project{
				Id:             row.ID,
				Name:           row.Name,
				Environments:   []*domain.Environment{},
				OrganizationId: row.OrganizationID,
			}
			projectsMap[proj.Id] = proj
		}

		proj.Environments = append(proj.Environments, &domain.Environment{
			Id:   row.EnvironmentID,
			Slug: row.EnvironmentSlug,
			Name: row.EnvironmentName,
		})
	}

	projects := make([]*domain.Project, 0, len(projectsMap))
	for _, p := range projectsMap {
		projects = append(projects, p)
	}

	return projects, nil
}

func (or *OrganizationRepository) GetOrganizationClusters(ctx context.Context, organizationID int64) ([]*domain.Cluster, error) {
	rows, err := or.queries.GetOrganizationClusters(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if rows == nil {
		return []*domain.Cluster{}, nil
	}

	clusters := make([]*domain.Cluster, 0, len(rows))
	for _, c := range rows {
		clusters = append(clusters, &domain.Cluster{
			Id:             c.ID,
			Name:           c.Name,
			OrganizationId: c.OrganizationID,
		})
	}

	return clusters, nil
}
