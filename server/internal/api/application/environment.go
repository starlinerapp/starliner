package application

import (
	"context"
	"fmt"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/core/domain/port"
	coreService "starliner.app/internal/core/domain/service"
)

type EnvironmentApplication struct {
	crypto                port.Crypto
	organizationService   *service.OrganizationService
	normalizerService     *coreService.NormalizerService
	environmentRepository interfaces.EnvironmentRepository
	projectRepository     interfaces.ProjectRepository
}

func NewEnvironmentApplication(
	crypto port.Crypto,
	normalizerService *coreService.NormalizerService,
	organizationService *service.OrganizationService,
	environmentRepository interfaces.EnvironmentRepository,
	projectRepository interfaces.ProjectRepository,
) *EnvironmentApplication {
	return &EnvironmentApplication{
		crypto:                crypto,
		normalizerService:     normalizerService,
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

	environmentSlug, err := ea.normalizerService.FormatToDNS1123(name)
	if err != nil {
		return err
	}

	project, err := ea.projectRepository.GetProject(ctx, userId, projectId)
	if err != nil {
		return err
	}

	namespace, err := ea.normalizerService.FormatToDNS1123(project.Name + "-" + name)
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

	git, err := ea.environmentRepository.GetEnvironmentGitDeployments(ctx, environmentId, userId)
	if err != nil {
		return nil, err
	}

	gitDeployments := make([]*value.GitDeployment, len(git))
	for i, d := range git {
		normalizedServiceName, err := ea.normalizerService.FormatToDNS1123(d.Name)
		if err != nil {
			return nil, err
		}

		internalEndpoint := fmt.Sprintf("%s:%s", normalizedServiceName, d.Port)
		gitDeployments[i] = value.NewGitDeployment(d, internalEndpoint)
	}

	images, err := ea.environmentRepository.GetEnvironmentImageDeployments(ctx, environmentId, userId)
	if err != nil {
		return nil, err
	}

	imageDeployments := make([]*value.ImageDeployment, len(images))
	for i, d := range images {
		normalizedServiceName, err := ea.normalizerService.FormatToDNS1123(d.ServiceName)
		if err != nil {
			return nil, err
		}

		internalEndpoint := fmt.Sprintf("%s:%s", normalizedServiceName, d.Port)
		imageDeployments[i] = value.NewImageDeployment(d, internalEndpoint)
	}

	databases, err := ea.environmentRepository.GetEnvironmentDatabaseDeployments(ctx, environmentId, userId)
	if err != nil {
		return nil, err
	}

	databaseDeployments := make([]*value.DatabaseDeployment, len(databases))
	for i, d := range databases {
		var password *string

		if d.Password != nil {
			decrypted, err := ea.crypto.Decrypt(*d.Password)
			if err != nil {
				return nil, err
			}
			password = &decrypted
		}

		normalizedServiceName, err := ea.normalizerService.FormatToDNS1123(d.ServiceName)
		if err != nil {
			return nil, err
		}

		internalEndpoint := fmt.Sprintf("%s-rw:%s", normalizedServiceName, d.Port)
		databaseDeployments[i] = &value.DatabaseDeployment{
			Id:               d.Id,
			ServiceName:      d.ServiceName,
			InternalEndpoint: internalEndpoint,
			Status:           d.Status,
			Database:         d.Database,
			Username:         d.Username,
			Password:         password,
			Port:             d.Port,
		}
	}

	return &value.Deployments{
		Databases:      databaseDeployments,
		Images:         imageDeployments,
		Ingresses:      value.NewIngressDeployments(ingresses),
		GitDeployments: gitDeployments,
	}, nil
}
