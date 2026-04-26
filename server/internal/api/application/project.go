package application

import (
	"context"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	coreService "starliner.app/internal/core/domain/service"
)

type ProjectApplication struct {
	normalizerService     *coreService.NormalizerService
	organizationService   *service.OrganizationService
	teamService           *service.TeamService
	projectRepository     interfaces.ProjectRepository
	environmentRepository interfaces.EnvironmentRepository
	environmentService    *service.EnvironmentService
}

func NewProjectApplication(
	normalizerService *coreService.NormalizerService,
	organizationService *service.OrganizationService,
	teamService *service.TeamService,
	projectRepository interfaces.ProjectRepository,
	environmentRepository interfaces.EnvironmentRepository,
	environmentService *service.EnvironmentService,
) *ProjectApplication {
	return &ProjectApplication{
		normalizerService:     normalizerService,
		organizationService:   organizationService,
		teamService:           teamService,
		projectRepository:     projectRepository,
		environmentRepository: environmentRepository,
		environmentService:    environmentService,
	}
}

func (pa *ProjectApplication) CreateProject(ctx context.Context, name string, clusterId int64, userId int64, teamId int64) (*value.Project, error) {
	err := pa.teamService.ValidateUserAndClusterInTeam(ctx, userId, teamId, clusterId)
	if err != nil {
		return nil, err
	}

	namespace, err := pa.normalizerService.FormatToDNS1123(name + "-" + value.EnvironmentProductionName)
	if err != nil {
		return nil, err
	}

	project, err := pa.projectRepository.CreateProjectWithEnvironment(
		ctx,
		name,
		namespace,
		value.EnvironmentProductionName,
		value.EnvironmentProductionSlug,
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
	envs, err := pa.projectRepository.GetProjectEnvironments(ctx, projectId, userId)
	if err != nil {
		return err
	}
	for _, env := range envs {
		if err := pa.environmentService.TearDownEnvironmentDeployments(ctx, env); err != nil {
			return err
		}
		if err := pa.environmentRepository.DeleteEnvironment(ctx, env.Id); err != nil {
			return err
		}
	}
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
