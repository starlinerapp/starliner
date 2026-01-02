package service

import (
	"context"
	"errors"
	"starliner.app/internal/domain"
	interfaces "starliner.app/internal/repository/interface"
	"strings"
)

type OrganizationService struct {
	organizationRepository interfaces.OrganizationRepository
}

func NewOrganizationService(organizationRepository interfaces.OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		organizationRepository: organizationRepository,
	}
}

func (os *OrganizationService) CreateOrganization(ctx context.Context, name string, ownerID int64) (*domain.Organization, error) {
	trimmed := strings.TrimSpace(name)
	organizationSlug := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")

	organization, err := os.organizationRepository.CreateOrganization(ctx, name, organizationSlug, ownerID)
	if err != nil {
		return nil, err
	}

	return organization, nil
}

func (os *OrganizationService) GetUserOrganizations(ctx context.Context, userID int64) ([]domain.Organization, error) {
	return os.organizationRepository.GetUserOrganizations(ctx, userID)
}

func (os *OrganizationService) GetProjectsForUser(ctx context.Context, userID int64, organizationID int64) ([]domain.Project, error) {
	err := os.ValidateUserOrganization(ctx, userID, organizationID)
	if err != nil {
		return nil, err
	}

	return os.organizationRepository.GetOrganizationProjects(ctx, organizationID)
}

func (os *OrganizationService) GetClustersForUser(ctx context.Context, userID int64, organizationID int64) ([]domain.Cluster, error) {
	err := os.ValidateUserOrganization(ctx, organizationID, userID)
	if err != nil {
		return nil, err
	}

	return os.organizationRepository.GetOrganizationClusters(ctx, organizationID)
}

func (os *OrganizationService) ValidateUserOrganization(ctx context.Context, organizationId int64, userId int64) error {
	organizations, err := os.organizationRepository.GetUserOrganizations(ctx, userId)
	if err != nil {
		return nil
	}

	found := false
	for _, org := range organizations {
		if org.Id == organizationId {
			found = true
			break
		}
	}
	if !found {
		return errors.New("organization not found")
	}

	return nil
}
