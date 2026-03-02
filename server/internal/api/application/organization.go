package application

import (
	"context"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/core/domain/port"
	"strings"
)

type OrganizationApplication struct {
	crypto                 port.Crypto
	organizationRepository interfaces.OrganizationRepository
	organizationService    *service.OrganizationService
}

func NewOrganizationApplication(
	crypto port.Crypto,
	organizationRepository interfaces.OrganizationRepository,
	organizationService *service.OrganizationService,
) *OrganizationApplication {
	return &OrganizationApplication{
		crypto:                 crypto,
		organizationRepository: organizationRepository,
		organizationService:    organizationService,
	}
}

func (oa *OrganizationApplication) CreateOrganization(ctx context.Context, name string, ownerID int64) error {
	trimmed := strings.TrimSpace(name)
	organizationSlug := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")

	_, err := oa.organizationRepository.CreateOrganization(ctx, name, organizationSlug, ownerID)
	if err != nil {
		return err
	}
	return nil
}

func (oa *OrganizationApplication) GetUserOrganizations(ctx context.Context, userID int64) ([]*value.Organization, error) {
	organizations, err := oa.organizationRepository.GetUserOrganizations(ctx, userID)
	if err != nil {
		return nil, err
	}
	return value.NewOrganizations(organizations), nil
}

func (oa *OrganizationApplication) GetProjectsForUser(ctx context.Context, userID int64, organizationID int64) ([]*value.Project, error) {
	err := oa.organizationService.ValidateUserInOrg(ctx, userID, organizationID)
	if err != nil {
		return nil, err
	}

	projects, err := oa.organizationRepository.GetOrganizationProjects(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	return value.NewProjects(projects), nil
}

func (oa *OrganizationApplication) GetClustersForUser(ctx context.Context, userID int64, organizationID int64) ([]*value.Cluster, error) {
	err := oa.organizationService.ValidateUserInOrg(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}

	clusters, err := oa.organizationRepository.GetOrganizationClusters(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	return value.NewClusters(clusters), nil
}

func (oa *OrganizationApplication) UpsertHetznerCredential(ctx context.Context, userID int64, organizationID int64, apiKey string) error {
	err := oa.organizationService.ValidateUserInOrg(ctx, organizationID, userID)
	if err != nil {
		return err
	}

	apiKeyEncrypted, err := oa.crypto.Encrypt(apiKey)
	if err != nil {
		return err
	}

	err = oa.organizationRepository.UpsertProvisioningCredentials(ctx, organizationID, apiKeyEncrypted, value.HetznerCredential)
	return err
}

func (oa *OrganizationApplication) GetHetznerCredential(ctx context.Context, userID int64, organizationID int64) (*value.ProvisioningCredential, error) {
	err := oa.organizationService.ValidateUserInOrg(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}

	c, err := oa.organizationRepository.GetOrganizationProvisioningCredential(ctx, organizationID, value.HetznerCredential)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	decrypted, err := oa.crypto.Decrypt(c.Secret)
	if err != nil {
		return nil, err
	}

	return &value.ProvisioningCredential{
		Provider: c.Provider,
		Secret:   decrypted,
	}, nil
}
