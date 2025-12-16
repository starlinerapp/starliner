package service

import (
	"context"
	"log"
	"starliner.app/pkg/domain"
	v1 "starliner.app/pkg/proto/v1"
	"starliner.app/pkg/queue"
	interfaces "starliner.app/pkg/repository/interface"
	"strings"
)

type ProjectService struct {
	organizationService    *OrganizationService
	projectRepository      interfaces.ProjectRepository
	organizationRepository interfaces.OrganizationRepository
	environmentRepository  interfaces.EnvironmentRepository
	projectPublisher       *queue.Publisher[*v1.Project]
}

func NewProjectService(
	organizationService *OrganizationService,
	projectRepository interfaces.ProjectRepository,
	organizationRepository interfaces.OrganizationRepository,
	environmentRepository interfaces.EnvironmentRepository,
	projectPublisher *queue.Publisher[*v1.Project],
) *ProjectService {
	return &ProjectService{
		organizationService:    organizationService,
		projectRepository:      projectRepository,
		organizationRepository: organizationRepository,
		environmentRepository:  environmentRepository,
		projectPublisher:       projectPublisher,
	}
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

	err = ps.projectPublisher.Publish("project.created", &v1.Project{
		Id:           project.Id,
		Name:         project.Name,
		Organization: project.OrganizationId,
	})
	log.Printf("error publishing: %v", err)

	environments := []domain.Environment{*environment}
	project.Environments = environments

	return project, nil
}

func (ps *ProjectService) GetProject(ctx context.Context, projectId int64, userId int64) (*domain.Project, error) {
	return ps.projectRepository.GetProject(ctx, projectId, userId)
}
