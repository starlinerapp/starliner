package service

import (
	"context"
	"starliner.app/pkg/domain"
	"starliner.app/pkg/repository"
	"strings"
)

type EnvironmentService struct {
	organizationService   *OrganizationService
	environmentRepository *repository.EnvironmentRepository
}

func NewEnvironmentService(environmentRepository *repository.EnvironmentRepository, organizationService *OrganizationService) *EnvironmentService {
	return &EnvironmentService{environmentRepository: environmentRepository, organizationService: organizationService}
}

func (es *EnvironmentService) CreateEnvironment(ctx context.Context, name string, userId int64, organizationId int64, projectId int64) (*domain.Environment, error) {
	err := es.organizationService.ValidateUserOrganization(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(name)
	environmentSlug := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")

	environment, err := es.environmentRepository.CreateEnvironment(ctx, name, environmentSlug, projectId)
	if err != nil {
		return nil, err
	}

	return environment, nil
}
