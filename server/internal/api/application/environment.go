package application

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/core/domain/port"
)

type EnvironmentApplication struct {
	crypto                port.Crypto
	organizationService   *service.OrganizationService
	namespaceService      *service.NamespaceService
	environmentRepository interfaces.EnvironmentRepository
	projectRepository     interfaces.ProjectRepository
}

func NewEnvironmentApplication(
	crypto port.Crypto,
	namespaceService *service.NamespaceService,
	organizationService *service.OrganizationService,
	environmentRepository interfaces.EnvironmentRepository,
	projectRepository interfaces.ProjectRepository,
) *EnvironmentApplication {
	return &EnvironmentApplication{
		crypto:                crypto,
		namespaceService:      namespaceService,
		organizationService:   organizationService,
		environmentRepository: environmentRepository,
		projectRepository:     projectRepository,
	}
}

func (ea *EnvironmentApplication) CreateEnvironment(
	ctx context.Context,
	name string,
	userId int64,
	organizationId int64,
	projectId int64,
) error {
	err := ea.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return err
	}

	environmentSlug, err := ea.namespaceService.FormatToDNS1123(name)
	if err != nil {
		return err
	}

	project, err := ea.projectRepository.GetProject(ctx, userId, projectId)
	if err != nil {
		return err
	}

	namespace, err := ea.namespaceService.FormatToDNS1123(project.Name + "-" + name)
	if err != nil {
		return err
	}

	_, err = ea.environmentRepository.CreateEnvironment(ctx, name, namespace, environmentSlug, projectId)
	if err != nil {
		return err
	}
	return nil
}

func (ea *EnvironmentApplication) GetEnvironmentDeployments(ctx context.Context, environmentId int64, userId int64) (*value.Deployments, error) {
	ingresses, err := ea.environmentRepository.GetEnvironmentIngressDeployments(ctx, environmentId, userId)
	if err != nil {
		return nil, err
	}

	images, err := ea.environmentRepository.GetEnvironmentImageDeployments(ctx, environmentId, userId)
	if err != nil {
		return nil, err
	}

	databases, err := ea.environmentRepository.GetEnvironmentDatabaseDeployments(ctx, environmentId, userId)
	if err != nil {
		return nil, err
	}

	databaseDeployments := make([]*entity.DatabaseDeployment, len(databases))
	for i, d := range databases {
		var password *string

		if d.Password != nil {
			decrypted, err := ea.crypto.Decrypt(*d.Password)
			if err != nil {
				return nil, err
			}
			password = &decrypted
		}

		databaseDeployments[i] = &entity.DatabaseDeployment{
			Id:            d.Id,
			ServiceName:   d.ServiceName,
			Status:        d.Status,
			Database:      d.Database,
			Username:      d.Username,
			Password:      password,
			Port:          d.Port,
			EnvironmentId: d.EnvironmentId,
		}
	}

	return &value.Deployments{
		Databases: value.NewDatabaseDeployments(databaseDeployments),
		Images:    value.NewImageDeployments(images),
		Ingresses: value.NewIngressDeployments(ingresses),
	}, nil
}
