package interfaces

import (
	"context"
	"starliner.app/internal/domain/entity"
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, name string, slug string, ownerID int64) (*entity.Organization, error)
	GetOrganization(ctx context.Context, id int64) (*entity.Organization, error)
	GetUserOrganizations(ctx context.Context, userID int64) ([]*entity.Organization, error)
	GetOrganizationProjects(ctx context.Context, organizationID int64) ([]*entity.ProjectWithEnvironments, error)
	GetOrganizationClusters(ctx context.Context, organizationID int64) ([]*entity.Cluster, error)
}
