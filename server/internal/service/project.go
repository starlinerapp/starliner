package service

import (
	"context"
	"starliner.app/internal/domain"
	interfaces "starliner.app/internal/repository/interface"
	"starliner.app/internal/service/model"
	"strings"
)

type ProjectService struct {
	organizationService    *OrganizationService
	projectRepository      interfaces.ProjectRepository
	organizationRepository interfaces.OrganizationRepository
	environmentRepository  interfaces.EnvironmentRepository
}

func NewProjectService(
	organizationService *OrganizationService,
	projectRepository interfaces.ProjectRepository,
	organizationRepository interfaces.OrganizationRepository,
	environmentRepository interfaces.EnvironmentRepository,
) *ProjectService {
	return &ProjectService{
		organizationService:    organizationService,
		projectRepository:      projectRepository,
		organizationRepository: organizationRepository,
		environmentRepository:  environmentRepository,
	}
}

func (ps *ProjectService) CreateProject(ctx context.Context, name string, organizationId int64, clusterId int64, userId int64) (*model.Project, error) {
	err := ps.organizationService.ValidateUserOrganization(ctx, organizationId, userId)
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

	environmentsModel := model.NewEnvironments([]*domain.Environment{environment})
	projectModel := model.NewProject(project)
	projectModel.Environments = environmentsModel

	return projectModel, nil
}

func (ps *ProjectService) GetProject(ctx context.Context, projectId int64, userId int64) (*model.Project, error) {
	project, err := ps.projectRepository.GetProject(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}
	return model.NewProject(project), nil
}
