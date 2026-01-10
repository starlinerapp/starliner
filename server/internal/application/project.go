package application

import (
	"context"
	"starliner.app/internal/domain/entity"
	"starliner.app/internal/domain/port"
	interfaces "starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/domain/service"
	"starliner.app/internal/domain/value"
	"strings"
)

type ProjectApplication struct {
	organizationService    *service.OrganizationService
	projectRepository      interfaces.ProjectRepository
	clusterRepository      interfaces.ClusterRepository
	organizationRepository interfaces.OrganizationRepository
	environmentRepository  interfaces.EnvironmentRepository
	deploy                 port.Deploy
	crypto                 port.Crypto
	queue                  port.Queue
}

func NewProjectApplication(
	organizationService *service.OrganizationService,
	projectRepository interfaces.ProjectRepository,
	organizationRepository interfaces.OrganizationRepository,
	clusterRepository interfaces.ClusterRepository,
	environmentRepository interfaces.EnvironmentRepository,
	deploy port.Deploy,
	crypto port.Crypto,
	queue port.Queue,
) *ProjectApplication {
	return &ProjectApplication{
		organizationService:    organizationService,
		projectRepository:      projectRepository,
		organizationRepository: organizationRepository,
		clusterRepository:      clusterRepository,
		environmentRepository:  environmentRepository,
		deploy:                 deploy,
		crypto:                 crypto,
		queue:                  queue,
	}
}

func (ps *ProjectApplication) CreateProject(ctx context.Context, name string, organizationId int64, clusterId int64, userId int64) (*value.Project, error) {
	err := ps.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	project, err := ps.projectRepository.CreateProject(ctx, name, organizationId, clusterId)
	if err != nil {
		return nil, err
	}

	productionEnvName := "Production"
	environment, err := ps.environmentRepository.CreateEnvironment(ctx, productionEnvName, strings.ToLower(productionEnvName), project.Id)
	if err != nil {
		return nil, err
	}

	environmentsModel := value.NewEnvironments([]*entity.Environment{environment})
	projectModel := value.NewProject(project)
	projectModel.Environments = environmentsModel

	return projectModel, nil
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
