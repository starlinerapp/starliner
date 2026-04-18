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
	environmentService    *service.EnvironmentService
	normalizerService     *coreService.NormalizerService
	environmentRepository interfaces.EnvironmentRepository
	projectRepository     interfaces.ProjectRepository
	buildRepository       interfaces.BuildRepository
}

func NewEnvironmentApplication(
	crypto port.Crypto,
	normalizerService *coreService.NormalizerService,
	organizationService *service.OrganizationService,
	environmentService *service.EnvironmentService,
	environmentRepository interfaces.EnvironmentRepository,
	projectRepository interfaces.ProjectRepository,
	buildRepository interfaces.BuildRepository,
) *EnvironmentApplication {
	return &EnvironmentApplication{
		crypto:                crypto,
		normalizerService:     normalizerService,
		organizationService:   organizationService,
		environmentService:    environmentService,
		environmentRepository: environmentRepository,
		projectRepository:     projectRepository,
		buildRepository:       buildRepository,
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

		latestArgs, err := ea.buildRepository.GetLatestBuildArgs(ctx, d.Id)
		var args []*value.Arg
		if err == nil {
			args = make([]*value.Arg, len(latestArgs))
			for j, a := range latestArgs {
				args[j] = &value.Arg{
					Name:  a.Name,
					Value: a.Value,
				}
			}
		} else {
			args = []*value.Arg{}
		}

		gitDeployments[i] = value.NewGitDeployment(d, internalEndpoint, args)
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

func (ea *EnvironmentApplication) GetEnvironmentGitDeploymentBuilds(ctx context.Context, userId int64, environmentId int64) ([]*value.GitDeploymentBuild, error) {
	err := ea.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return nil, err
	}

	builds, err := ea.environmentRepository.GetEnvironmentGitDeploymentBuilds(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	valueBuilds := make([]*value.GitDeploymentBuild, len(builds))
	for i, b := range builds {
		args := make([]*value.Arg, len(b.Args))
		for j, a := range b.Args {
			args[j] = &value.Arg{
				Name:  a.Name,
				Value: a.Value,
			}
		}

		valueBuilds[i] = &value.GitDeploymentBuild{
			BuildId:        b.BuildId,
			DeploymentId:   b.DeploymentId,
			DeploymentName: b.DeploymentName,
			CommitHash:     b.CommitHash,
			Source:         b.Source,
			Status:         value.BuildStatus(b.Status),
			GitUrl:         b.GitUrl,
			ProjectPath:    b.ProjectPath,
			DockerfilePath: b.DockerfilePath,
			CreatedAt:      b.CreatedAt,
			Args:           args,
		}
	}
	return valueBuilds, nil
}
