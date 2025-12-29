package interfaces

import (
	"context"
	"starliner.app/pkg/domain"
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, name string, slug string, ownerID int64) (*domain.Organization, error)
	GetUserOrganizations(ctx context.Context, userID int64) ([]domain.Organization, error)
	GetOrganizationProjects(ctx context.Context, organizationID int64) ([]domain.Project, error)
	GetOrganizationClusters(ctx context.Context, organizationID int64) ([]domain.Cluster, error)
}
