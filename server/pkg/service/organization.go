package service

import (
	"context"
	"errors"
	"starliner.app/pkg/domain"
	"starliner.app/pkg/repository"
	"strings"
)

type OrganizationService struct {
	organizationRepository *repository.OrganizationRepository
}

func NewOrganizationService(organizationRepository *repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		organizationRepository: organizationRepository,
	}
}

func (os *OrganizationService) CreateOrganization(ctx context.Context, name string, ownerID int64) (*domain.Organization, error) {
	trimmed := strings.TrimSpace(name)
	slug := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")

	return os.organizationRepository.CreateOrganization(ctx, name, slug, ownerID)
}

func (os *OrganizationService) GetUserOrganizations(ctx context.Context, userID int64) ([]domain.Organization, error) {
	return os.organizationRepository.GetUserOrganizations(ctx, userID)
}

func (os *OrganizationService) GetProjectsForUser(ctx context.Context, userID int64, organizationID int64) ([]domain.Project, error) {
	organizations, err := os.organizationRepository.GetUserOrganizations(ctx, userID)
	if err != nil {
		return nil, err
	}

	found := false
	for _, org := range organizations {
		if org.Id == organizationID {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("organization not found")
	}

	return os.organizationRepository.GetOrganizationProjects(ctx, organizationID)
}
