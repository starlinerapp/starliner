package application

import (
	"context"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	"strings"
)

type OrganizationApplication struct {
	organizationRepository interfaces.OrganizationRepository
	organizationService    *service.OrganizationService
}

func NewOrganizationApplication(
	organizationRepository interfaces.OrganizationRepository,
	organizationService *service.OrganizationService,
) *OrganizationApplication {
	return &OrganizationApplication{
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
