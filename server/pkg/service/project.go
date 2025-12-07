package service

import (
	"context"
	"starliner.app/pkg/domain"
	interfaces "starliner.app/pkg/repository/interface"
	"strings"
)

type ProjectService struct {
	organizationService    *OrganizationService
	projectRepository      interfaces.ProjectRepository
	organizationRepository interfaces.OrganizationRepository
	environmentRepository  interfaces.EnvironmentRepository
}

func NewProjectService(organizationService *OrganizationService, projectRepository interfaces.ProjectRepository, organizationRepository interfaces.OrganizationRepository, environmentRepository interfaces.EnvironmentRepository) *ProjectService {
	return &ProjectService{organizationService: organizationService, projectRepository: projectRepository, organizationRepository: organizationRepository, environmentRepository: environmentRepository}
}

func (ps *ProjectService) CreateProject(ctx context.Context, name string, organizationId int64, userId int64) (*domain.Project, error) {
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

	environments := []domain.Environment{*environment}
	project.Environments = environments

	return project, nil
}

func (ps *ProjectService) GetProject(ctx context.Context, projectId int64, userId int64) (*domain.Project, error) {
	return ps.projectRepository.GetProject(ctx, projectId, userId)
}
