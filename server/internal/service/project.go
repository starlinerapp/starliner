package service

import (
	"context"
	"log"
	"starliner.app/internal/domain"
	"starliner.app/internal/infrastructure/queue"
	v1 "starliner.app/internal/infrastructure/queue/proto/v1"
	interfaces "starliner.app/internal/repository/interface"
	"starliner.app/internal/service/model"
	"strconv"
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

	err = ps.projectPublisher.Publish(queue.CreateProject, strconv.FormatInt(project.Id, 10), &v1.Project{
		Id:             project.Id,
		Name:           project.Name,
		OrganizationId: project.OrganizationId,
		ClusterId:      *project.ClusterId,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
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

func (ps *ProjectService) DeleteProject(ctx context.Context, projectId int64, userId int64) error {
	return ps.projectRepository.DeleteProject(ctx, projectId, userId)
}
