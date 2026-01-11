package application

import (
	"context"
	"starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/domain/service"
	"strings"
)

type EnvironmentApplication struct {
	organizationService   *service.OrganizationService
	environmentRepository interfaces.EnvironmentRepository
}

func NewEnvironmentApplication(
	environmentRepository interfaces.EnvironmentRepository,
	organizationService *service.OrganizationService,
) *EnvironmentApplication {
	return &EnvironmentApplication{
		environmentRepository: environmentRepository,
		organizationService:   organizationService,
	}
}

func (ea *EnvironmentApplication) CreateEnvironment(
	ctx context.Context,
	name string,
	userId int64,
	organizationId int64,
	projectId int64,
) error {
	err := ea.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return err
	}

	trimmed := strings.TrimSpace(name)
	environmentSlug := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")

	_, err = ea.environmentRepository.CreateEnvironment(ctx, name, environmentSlug, projectId)
	if err != nil {
		return err
	}
	return nil
}
