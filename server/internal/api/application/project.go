package application

import (
	"context"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	"strings"
)

type ProjectApplication struct {
	namespaceService      *service.NamespaceService
	organizationService   *service.OrganizationService
	projectRepository     interfaces.ProjectRepository
	environmentRepository interfaces.EnvironmentRepository
}

func NewProjectApplication(
	namespaceService *service.NamespaceService,
	organizationService *service.OrganizationService,
	projectRepository interfaces.ProjectRepository,
	environmentRepository interfaces.EnvironmentRepository,
) *ProjectApplication {
	return &ProjectApplication{
		namespaceService:      namespaceService,
		organizationService:   organizationService,
		projectRepository:     projectRepository,
		environmentRepository: environmentRepository,
	}
}

func (ps *ProjectApplication) CreateProject(ctx context.Context, name string, organizationId int64, clusterId int64, userId int64) (*value.Project, error) {
	err := ps.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	productionEnvName := "Production"
	namespace, err := ps.namespaceService.FormatToDNS1123(name + "-" + productionEnvName)
	if err != nil {
		return nil, err
	}

	project, err := ps.projectRepository.CreateProjectWithEnvironment(
		ctx,
		name,
		namespace,
		productionEnvName,
		strings.ToLower(productionEnvName),
		organizationId,
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
