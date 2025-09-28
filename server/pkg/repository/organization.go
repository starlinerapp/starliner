package repository

import (
	"context"
	"starliner.app/pkg/db/sqlc"
	"starliner.app/pkg/domain"
)

type OrganizationRepository struct {
	queries *sqlc.Queries
}

func NewOrganizationRepository(queries *sqlc.Queries) *OrganizationRepository {
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

func (or *OrganizationRepository) GetUserOrganizations(ctx context.Context, userID int64) ([]domain.Organization, error) {
	organizations, err := or.queries.GetUserOrganizations(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Organization, len(organizations))
	for i, organization := range organizations {
		result[i] = domain.Organization{
			Id:      organization.ID,
			Name:    organization.Name,
			Slug:    organization.Slug,
			OwnerId: organization.OwnerID,
		}
	}

	return result, nil
}

func (or *OrganizationRepository) GetOrganizationProjects(ctx context.Context, organizationID int64) ([]domain.Project, error) {
	projects, err := or.queries.GetOrganizationProjects(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Project, len(projects))
	for i, project := range projects {
		result[i] = domain.Project{
			Id:             project.ID,
			Name:           project.Name,
			OrganizationId: project.OrganizationID,
		}
	}

	return result, nil
}
