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
	teamService           *service.TeamService
	projectRepository     interfaces.ProjectRepository
	environmentRepository interfaces.EnvironmentRepository
}

func NewProjectApplication(
	normalizerService *coreService.NormalizerService,
	organizationService *service.OrganizationService,
	teamService *service.TeamService,
	projectRepository interfaces.ProjectRepository,
	environmentRepository interfaces.EnvironmentRepository,
) *ProjectApplication {
	return &ProjectApplication{
		normalizerService:     normalizerService,
		organizationService:   organizationService,
		teamService:           teamService,
		projectRepository:     projectRepository,
		environmentRepository: environmentRepository,
	}
}

func (ps *ProjectApplication) CreateProject(ctx context.Context, name string, organizationId int64, clusterId int64, userId int64, teamId int64) (*value.Project, error) {
	err := ps.teamService.ValidateUserAndClusterInTeam(ctx, userId, teamId, clusterId)
	if err != nil {
		return nil, err
	}

	productionEnvName := "Production"
	namespace, err := ps.normalizerService.FormatToDNS1123(name + "-" + productionEnvName)
	if err != nil {
		return nil, err
	}

	project, err := ps.projectRepository.CreateProjectWithEnvironment(
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

func (ps *ProjectApplication) GetProject(ctx context.Context, projectId int64, userId int64) (*value.Project, error) {
	project, err := ps.projectRepository.GetProject(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}
	return value.NewProject(project), nil
}

func (ps *ProjectApplication) DeleteProject(ctx context.Context, projectId int64, userId int64) error {
	return ps.projectRepository.DeleteProject(ctx, projectId, userId)
}

func (ps *ProjectApplication) GetProjectCluster(ctx context.Context, projectId int64, userId int64) (*value.ProjectCluster, error) {

	cluster, err := ps.projectRepository.GetProjectCluster(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}
	return value.NewProjectCluster(cluster), nil
}

func (ps *ProjectApplication) GetProjectEnvironments(ctx context.Context, projectId int64, userId int64) ([]*value.Environment, error) {
	environments, err := ps.projectRepository.GetProjectEnvironments(ctx, projectId, userId)
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
