package interfaces

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/value"
	"time"
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, name string, slug string, ownerID int64) (*entity.Organization, error)
	GetOrganization(ctx context.Context, id int64) (*entity.Organization, error)
	GetUserOrganizations(ctx context.Context, userID int64) ([]*entity.Organization, error)
	GetUserProjects(ctx context.Context, organizationID int64, userID int64) ([]*entity.Project, error)
	GetOrganizationClusters(ctx context.Context, organizationID int64) ([]*entity.Cluster, error)
	UpsertProvisioningCredentials(ctx context.Context, organizationID int64, apiKey string, provider value.ProvisioningCredentialProvider) error
	GetOrganizationProvisioningCredential(
		ctx context.Context,
		organizationID int64,
		provider value.ProvisioningCredentialProvider,
	) (*value.ProvisioningCredential, error)
	AddOrganizationMember(ctx context.Context, organizationID int64, userID int64) error
	RemoveOrganizationMember(ctx context.Context, organizationID int64, userID int64) error
	CreateOrganizationInvite(ctx context.Context, organizationID int64, expiresAt time.Time) (*entity.OrganizationInvite, error)
	GetOrganizationInviteById(ctx context.Context, inviteId string) (*entity.OrganizationInvite, error)
	GetOrganizationMembers(ctx context.Context, organizationID int64) ([]*entity.User, error)
}
