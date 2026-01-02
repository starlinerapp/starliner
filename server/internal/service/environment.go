package service

import (
	"context"
	interfaces "starliner.app/internal/repository/interface"
	"strings"
)

type EnvironmentService struct {
	organizationService   *OrganizationService
	environmentRepository interfaces.EnvironmentRepository
}

func NewEnvironmentService(environmentRepository interfaces.EnvironmentRepository, organizationService *OrganizationService) *EnvironmentService {
	return &EnvironmentService{environmentRepository: environmentRepository, organizationService: organizationService}
}

func (es *EnvironmentService) CreateEnvironment(ctx context.Context, name string, userId int64, organizationId int64, projectId int64) error {
	err := es.organizationService.ValidateUserOrganization(ctx, organizationId, userId)
	if err != nil {
		return err
	}

	trimmed := strings.TrimSpace(name)
	environmentSlug := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")

	_, err = es.environmentRepository.CreateEnvironment(ctx, name, environmentSlug, projectId)
	if err != nil {
		return err
	}
	return nil
}
