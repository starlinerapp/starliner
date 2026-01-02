package service

import (
	"context"
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

func (ps *ProjectService) CreateProject(ctx context.Context, name string, organizationId int64, userId int64) (*model.Project, error) {
	err := ps.organizationService.ValidateUserOrganization(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	project, err := ps.projectRepository.CreateProject(ctx, name, organizationId)
	if err != nil {
		return nil, err
	}

	productionEnvName := "Production"
	environment, err := ps.environmentRepository.CreateEnvironment(ctx, productionEnvName, strings.ToLower(productionEnvName), project.Id)
	if err != nil {
		return nil, err
	}

	environments := []model.Environment{
		{
			Id:   environment.Id,
			Slug: environment.Slug,
			Name: environment.Name,
		},
	}

	projectModel := &model.Project{
		Id:             project.Id,
		Name:           project.Name,
		Environments:   environments,
		OrganizationId: project.OrganizationId,
	}
	projectModel.Environments = environments

	return projectModel, nil
}

func (ps *ProjectService) GetProject(ctx context.Context, projectId int64, userId int64) (*model.Project, error) {
	project, err := ps.projectRepository.GetProject(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}

	environmentModels := make([]model.Environment, len(project.Environments))
	for i, env := range project.Environments {
		environmentModels[i] = model.Environment{
			Id:   env.Id,
			Slug: env.Slug,
			Name: env.Name,
		}
	}

	return &model.Project{
		Id:             project.Id,
		Name:           project.Name,
		Environments:   environmentModels,
		OrganizationId: project.OrganizationId,
	}, nil
}
