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

func (ea *EnvironmentApplication) GetEnvironmentDeployments(ctx context.Context, environmentId int64, userId int64) (*value.Deployments, error) {
	images, err := ea.environmentRepository.GetEnvironmentImageDeployments(ctx, environmentId, userId)
	if err != nil {
		return nil, err
	}

	databases, err := ea.environmentRepository.GetEnvironmentDatabaseDeployments(ctx, environmentId, userId)
	if err != nil {
		return nil, err
	}

	return &value.Deployments{
		Databases: value.NewDatabaseDeployments(databases),
		Images:    value.NewImageDeployments(images),
	}, nil
}
