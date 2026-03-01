package interfaces

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/value"
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, name string, slug string, ownerID int64) (*entity.Organization, error)
	GetOrganization(ctx context.Context, id int64) (*entity.Organization, error)
	GetUserOrganizations(ctx context.Context, userID int64) ([]*entity.Organization, error)
	GetOrganizationProjects(ctx context.Context, organizationID int64) ([]*entity.Project, error)
	GetOrganizationClusters(ctx context.Context, organizationID int64) ([]*entity.Cluster, error)
	UpsertProvisioningCredentials(ctx context.Context, organizationID int64, apiKey string, provider value.CredentialProvider) error
}
