package application

import (
	"context"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	"strings"
)

type EnvironmentApplication struct {
	organizationService   *service.OrganizationService
	environmentRepository _interface.EnvironmentRepository
}

func NewEnvironmentApplication(
	environmentRepository _interface.EnvironmentRepository,
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

func (ea *EnvironmentApplication) GetEnvironmentDeployments(ctx context.Context, environmentId int64, userId int64) ([]*value.Deployment, error) {
	deployments, err := ea.environmentRepository.GetEnvironmentDeployments(ctx, environmentId, userId)
	if err != nil {
		return nil, err
	}

	return value.NewDeployments(deployments), nil
}
