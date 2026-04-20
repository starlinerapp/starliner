package application

import (
	"context"
	"fmt"
	"log"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	corePort "starliner.app/internal/core/domain/port"
	coreService "starliner.app/internal/core/domain/service"
	coreValue "starliner.app/internal/core/domain/value"
	"strconv"
)

type EnvironmentApplication struct {
	cfg                   *conf.Config
	crypto                corePort.Crypto
	queue                 port.Queue
	organizationService   *service.OrganizationService
	environmentService    *service.EnvironmentService
	normalizerService     *coreService.NormalizerService
	parserService         *service.ParserService
	resolverService       *service.ResolverService
	buildRepository       interfaces.BuildRepository
	environmentRepository interfaces.EnvironmentRepository
	projectRepository     interfaces.ProjectRepository
}

func NewEnvironmentApplication(
	cfg *conf.Config,
	crypto corePort.Crypto,
	queue port.Queue,
	normalizerService *coreService.NormalizerService,
	parserService *service.ParserService,
	resolverService *service.ResolverService,
	organizationService *service.OrganizationService,
	environmentService *service.EnvironmentService,
	buildRepository interfaces.BuildRepository,
	environmentRepository interfaces.EnvironmentRepository,
	projectRepository interfaces.ProjectRepository,
) *EnvironmentApplication {
	return &EnvironmentApplication{
		cfg:                   cfg,
		crypto:                crypto,
		queue:                 queue,
		normalizerService:     normalizerService,
		parserService:         parserService,
		resolverService:       resolverService,
		organizationService:   organizationService,
		environmentService:    environmentService,
		buildRepository:       buildRepository,
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
	sourceEnvironmentId *int64,
) (*value.Environment, error) {
	err := ea.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	environmentSlug, err := ea.normalizerService.FormatToDNS1123(name)
	if err != nil {
		return nil, err
	}

	project, err := ea.projectRepository.GetProject(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}

	namespace, err := ea.normalizerService.FormatToDNS1123(project.Name + "-" + name)
	if err != nil {
		return nil, err
	}

	if sourceEnvironmentId != nil {
		randomPrefix := ea.environmentService.RandomPrefix(4)
		env, err := ea.environmentRepository.DuplicateEnvironment(ctx, userId, name, namespace, environmentSlug, projectId, *sourceEnvironmentId, randomPrefix)
		if err != nil {
			return nil, err
		}
		deployments, err := ea.GetEnvironmentDeployments(ctx, env.Id, userId)
		if err != nil {
			return nil, err
		}

		cluster, err := ea.environmentRepository.GetEnvironmentCluster(ctx, env.Id)
		if err != nil {
			return nil, err
		}
		kubeconfigBase64, err := ea.crypto.Decrypt(*cluster.Kubeconfig)
		if err != nil {
			return nil, err
		}

		ingressDeployments := deployments.Ingresses
		for _, d := range ingressDeployments {
			coreHosts := make([]coreValue.IngressHost, 0, len(d.IngressHosts))
			for _, h := range d.IngressHosts {
				ch := coreValue.IngressHost{
					Host: h.Host,
				}
				ch.Paths = make([]coreValue.IngressPath, 0, len(h.Paths))

				for _, p := range h.Paths {
					target, err := ea.environmentRepository.GetEnvironmentDeploymentByName(ctx, p.ServiceName, env.Id)
					if err != nil {
						return nil, err
					}

					targetPort, err := strconv.Atoi(target.Port)
					if err != nil {
						return nil, err
					}

					normalizedServiceName, err := ea.normalizerService.FormatToDNS1123(p.ServiceName)
					if err != nil {
						return nil, err
					}

					ch.Paths = append(ch.Paths, coreValue.IngressPath{
						Path:        p.Path,
						PathType:    coreValue.PathType(p.PathType),
						ServiceName: normalizedServiceName,
						ServicePort: targetPort,
					})
				}
				coreHosts = append(coreHosts, ch)
			}
			err = ea.queue.PublishDeployIngress(&coreValue.IngressDeployment{
				IngressHosts:     coreHosts,
				DeploymentId:     d.Id,
				DeploymentName:   d.ServiceName,
				Namespace:        env.Namespace,
				KubeconfigBase64: kubeconfigBase64,
			})
			if err != nil {
				log.Printf("error publishing: %v", err)
			}
		}

		databaseDeployments := deployments.Databases
		for _, d := range databaseDeployments {
			normalizedServiceName, err := ea.normalizerService.FormatToDNS1123(d.ServiceName)
			if err != nil {
				return nil, err
			}

			err = ea.queue.PublishDeployDatabase(&coreValue.Deployment{
				Namespace:        env.Namespace,
				DeploymentId:     d.Id,
				DeploymentName:   normalizedServiceName,
				KubeconfigBase64: kubeconfigBase64,
			})
			if err != nil {
				log.Printf("error publishing: %v", err)
			}
		}

		imageDeployments := deployments.Images
		for _, d := range imageDeployments {
			deploymentPort, err := strconv.Atoi(d.Port)
			if err != nil {
				return nil, err
			}
			coreEnvs := value.ToCoreEnvVars(d.EnvVars)

			normalizedDeploymentName, err := ea.normalizerService.FormatToDNS1123(d.ServiceName)
			if err != nil {
				return nil, err
			}
			err = ea.queue.PublishDeployImage(&coreValue.ImageDeployment{
				DeploymentId:     d.Id,
				DeploymentName:   normalizedDeploymentName,
				Namespace:        env.Namespace,
				KubeconfigBase64: kubeconfigBase64,
				ImageName:        d.ImageName,
				ImageTag:         d.Tag,
				Port:             deploymentPort,
				VolumeSizeMiB:    d.VolumeSizeMiB,
				VolumeMountPath:  d.VolumeMountPath,
				EnvVars:          coreEnvs,
			})
			if err != nil {
				log.Printf("error publishing: %v", err)
			}
		}

		gitDeployments := deployments.GitDeployments
		for _, d := range gitDeployments {
			latestBuild, err := ea.buildRepository.GetLatestGitDeploymentBuild(ctx, *sourceEnvironmentId, d.ServiceName)
			if err != nil {
				return nil, err
			}

			coreEnvs := make([]*coreValue.EnvVar, 0, len(d.EnvVars))
			for _, e := range d.EnvVars {
				res, err := ea.parserService.Parse(e.Value)
				if err != nil {
					log.Printf("failed to parse env var: %v\n", err)
					continue
				}

				resolvedValue, err := ea.resolverService.Resolve(ctx, env.Id, res)
				if err != nil {
					log.Printf("failed to resolve env var: %v\n", err)
					continue
				}

				coreEnvs = append(coreEnvs, &coreValue.EnvVar{
					Name:  e.Name,
					Value: resolvedValue,
				})
			}

			normalizedDeploymentName, err := ea.normalizerService.FormatToDNS1123(d.ServiceName)
			if err != nil {
				return nil, err
			}
			deploymentPort, err := strconv.Atoi(d.Port)
			if err != nil {
				return nil, err
			}

			if latestBuild.ImageName == nil {
				return nil, fmt.Errorf("latest build for git deployment %s is nil", d.ServiceName)
			}

			err = ea.queue.PublishDeployImage(&coreValue.ImageDeployment{
				DeploymentId:     d.Id,
				DeploymentName:   normalizedDeploymentName,
				Namespace:        env.Namespace,
				KubeconfigBase64: kubeconfigBase64,
				ImageName:        *latestBuild.ImageName,
				ImageTag:         *latestBuild.CommitHash,
				Port:             deploymentPort,
				EnvVars:          coreEnvs,
			})
			if err != nil {
				log.Printf("failed to publish: %v\n", err)
			}
		}
		return value.NewEnvironment(env), nil
	}

	env, err := ea.environmentRepository.CreateEnvironment(ctx, name, namespace, environmentSlug, projectId)
	if err != nil {
		return nil, err
	}

	return value.NewEnvironment(env), nil
}

func (ea *EnvironmentApplication) GetEnvironmentDeployments(ctx context.Context, environmentId int64, userId int64) (*value.Deployments, error) {
	ingresses, err := ea.environmentRepository.GetUserEnvironmentIngressDeployments(ctx, environmentId, userId)
	if err != nil {
		return nil, err
	}

	git, err := ea.environmentRepository.GetUserEnvironmentGitDeployments(ctx, environmentId, userId)
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

	images, err := ea.environmentRepository.GetUserEnvironmentImageDeployments(ctx, environmentId, userId)
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

	databases, err := ea.environmentRepository.GetUserEnvironmentDatabaseDeployments(ctx, environmentId, userId)
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

func (ea *EnvironmentApplication) GetEnvironmentBranch(ctx context.Context, userId int64, environmentId int64) (string, error) {
	err := ea.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return "", err
	}

	branch, err := ea.environmentRepository.GetEnvironmentBranch(ctx, environmentId)
	if err != nil {
		return "", err
	}
	return branch, nil
}

func (ea *EnvironmentApplication) UpdateEnvironmentBranch(ctx context.Context, userId int64, environmentId int64, branch string) error {
	err := ea.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return err
	}

	err = ea.environmentRepository.UpdateEnvironmentBranch(ctx, environmentId, branch)
	if err != nil {
		return err
	}
	return nil
}
