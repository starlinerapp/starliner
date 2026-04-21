package application

import (
	"context"
	"strings"

	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	coreService "starliner.app/internal/core/domain/service"
)

type ProjectApplication struct {
	normalizerService     *coreService.NormalizerService
	organizationService   *service.OrganizationService
	projectRepository     interfaces.ProjectRepository
	environmentRepository interfaces.EnvironmentRepository
}

func NewProjectApplication(
	normalizerService *coreService.NormalizerService,
	organizationService *service.OrganizationService,
	projectRepository interfaces.ProjectRepository,
	environmentRepository interfaces.EnvironmentRepository,
) *ProjectApplication {
	return &ProjectApplication{
		normalizerService:     normalizerService,
		organizationService:   organizationService,
		projectRepository:     projectRepository,
		environmentRepository: environmentRepository,
	}
}

func (pa *ProjectApplication) CreateProject(ctx context.Context, name string, organizationId int64, clusterId int64, userId int64, teamId int64) (*value.Project, error) {
	err := pa.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	productionEnvName := "Production"
	namespace, err := pa.normalizerService.FormatToDNS1123(name + "-" + productionEnvName)
	if err != nil {
		return nil, err
	}

	project, err := pa.projectRepository.CreateProjectWithEnvironment(
		ctx,
		name,
		namespace,
		productionEnvName,
		strings.ToLower(productionEnvName),
		teamId,
		clusterId,
	)
	if err != nil {
		return nil, err
	}

	return value.NewProject(project), nil
}

func (pa *ProjectApplication) GetProject(ctx context.Context, projectId int64, userId int64) (*value.Project, error) {
	project, err := pa.projectRepository.GetProject(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}
	return value.NewProject(project), nil
}

func (pa *ProjectApplication) DeleteProject(ctx context.Context, projectId int64, userId int64) error {
	return pa.projectRepository.DeleteProject(ctx, projectId, userId)
}

func (pa *ProjectApplication) GetProjectCluster(ctx context.Context, projectId int64, userId int64) (*value.ProjectCluster, error) {

	cluster, err := pa.projectRepository.GetProjectCluster(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}
	return value.NewProjectCluster(cluster), nil
}

func (pa *ProjectApplication) GetProjectEnvironments(ctx context.Context, projectId int64, userId int64) ([]*value.Environment, error) {
	environments, err := pa.projectRepository.GetProjectEnvironments(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}

	valueEnvironments := make([]*value.Environment, len(environments))
	for i, e := range environments {
		valueEnvironments[i] = &value.Environment{
			Id:   e.Id,
			Name: e.Name,
			Slug: e.Slug,
		}
	}
	return valueEnvironments, nil
}

func (pa *ProjectApplication) GetProjectPreviewEnvironmentEnabled(ctx context.Context, projectId int64, userId int64) (bool, error) {
	enabled, err := pa.projectRepository.GetProjectPreviewEnvironmentEnabled(ctx, projectId, userId)
	if err != nil {
		return false, err
	}
	return enabled, nil
}

func (pa *ProjectApplication) ToggleProjectPreviewEnvironmentEnabled(ctx context.Context, projectId int64, userId int64) (bool, error) {
	enabled, err := pa.projectRepository.ToggleProjectPreviewEnvironmentEnabled(ctx, projectId, userId)
	if err != nil {
		return false, err
	}
	return enabled, nil
}
