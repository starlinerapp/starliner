package application

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"

	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	corePort "starliner.app/internal/core/domain/port"
	coreService "starliner.app/internal/core/domain/service"
	coreValue "starliner.app/internal/core/domain/value"
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
		env, err := ea.environmentRepository.DuplicateEnvironment(ctx, name, namespace, environmentSlug, projectId, *sourceEnvironmentId, randomPrefix, nil)
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
		if cluster.IPv4Address == nil || *cluster.IPv4Address == "" {
			return nil, fmt.Errorf("cluster ipv4 address is not set")
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
			err = createDeployOnlyBuild(ctx, ea.buildRepository, d.Id, value.BuildSourceDuplicate)
			if err != nil {
				log.Printf("failed to create ingress deploy build: %v", err)
				continue
			}
			err = ea.queue.PublishDeployIngress(&coreValue.IngressDeployment{
				IngressHosts:     coreHosts,
				DeploymentId:     d.Id,
				DeploymentName:   d.ServiceName,
				Namespace:        env.Namespace,
				KubeconfigBase64: kubeconfigBase64,
				ExpectedIP:       *cluster.IPv4Address,
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

			err = createDeployOnlyBuild(ctx, ea.buildRepository, d.Id, value.BuildSourceDuplicate)
			if err != nil {
				log.Printf("failed to create database deploy build: %v", err)
				continue
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
			err = createDeployOnlyBuild(ctx, ea.buildRepository, d.Id, value.BuildSourceDuplicate)
			if err != nil {
				log.Printf("failed to create image deploy build: %v", err)
				continue
			}
			err = ea.queue.PublishDeployImage(&coreValue.ImageDeployment{
				DeploymentId:          d.Id,
				DeploymentName:        normalizedDeploymentName,
				Namespace:             env.Namespace,
				KubeconfigBase64:      kubeconfigBase64,
				ImageRegistryUrl:      ea.cfg.ImageRegistryUrl,
				ImageRegistryUsername: ea.cfg.ImageRegistryUsername,
				ImageRegistryPassword: ea.cfg.ImageRegistryPassword,
				ImageName:             d.ImageName,
				ImageTag:              d.Tag,
				Port:                  deploymentPort,
				VolumeSizeMiB:         d.VolumeSizeMiB,
				VolumeMountPath:       d.VolumeMountPath,
				EnvVars:               coreEnvs,
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

			if latestBuild.ImageName == nil {
				return nil, fmt.Errorf("latest build for git deployment %s is nil", d.ServiceName)
			}

			err = ea.triggerDuplicateGitDeploy(ctx, d, latestBuild, env, kubeconfigBase64)
			if err != nil {
				return nil, err
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

		internalEndpoint := fmt.Sprintf("%s:%s", normalizedServiceName, d.Port)
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

	ingressBuilds, err := ea.environmentRepository.GetEnvironmentIngressDeploymentBuilds(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	imageBuilds, err := ea.environmentRepository.GetEnvironmentImageDeploymentBuilds(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	databaseBuilds, err := ea.environmentRepository.GetEnvironmentDatabaseDeploymentBuilds(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	allBuilds := append(builds, ingressBuilds...)
	allBuilds = append(allBuilds, imageBuilds...)
	allBuilds = append(allBuilds, databaseBuilds...)
	sort.Slice(allBuilds, func(i, j int) bool {
		return allBuilds[i].CreatedAt.After(allBuilds[j].CreatedAt)
	})

	valueBuilds := make([]*value.GitDeploymentBuild, len(allBuilds))
	for i, b := range allBuilds {
		valueBuilds[i] = &value.GitDeploymentBuild{
			BuildId:                 b.BuildId,
			DeploymentId:            b.DeploymentId,
			DeploymentName:          b.DeploymentName,
			DeploymentRolloutStatus: b.DeploymentRolloutStatus,
			CommitHash:              b.CommitHash,
			Source:                  b.Source,
			Status:                  value.BuildStatus(b.Status),
			CreatedAt:               b.CreatedAt,
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

func (ea *EnvironmentApplication) DeleteEnvironment(ctx context.Context, userId int64, environmentId int64) error {
	if err := ea.environmentService.ValidateUserPermission(ctx, userId, environmentId); err != nil {
		return err
	}
	env, err := ea.environmentRepository.GetEnvironmentById(ctx, environmentId)
	if err != nil {
		return err
	}

	if env.Slug == value.EnvironmentProductionSlug {
		return fmt.Errorf("cannot delete production environment")
	}

	if err := ea.environmentService.TearDownEnvironmentDeployments(ctx, env); err != nil {
		return err
	}
	return ea.environmentRepository.DeleteEnvironment(ctx, environmentId)
}

func (ea *EnvironmentApplication) triggerDuplicateGitDeploy(
	ctx context.Context,
	deployment *value.GitDeployment,
	sourceBuild *entity.GitDeploymentBuild,
	env *entity.Environment,
	kubeconfigBase64 string,
) error {
	coreEnvs := make([]*coreValue.EnvVar, 0, len(deployment.EnvVars))
	for _, e := range deployment.EnvVars {
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

	normalizedDeploymentName, err := ea.normalizerService.FormatToDNS1123(deployment.ServiceName)
	if err != nil {
		return err
	}
	deploymentPort, err := strconv.Atoi(deployment.Port)
	if err != nil {
		return err
	}

	b, err := ea.buildRepository.CreateBuild(ctx, deployment.Id, value.BuildSourceDuplicate)
	if err != nil {
		return err
	}

	err = ea.buildRepository.UpdateBuild(
		ctx,
		b.Id,
		value.BuildStatusSuccess,
		sourceBuild.CommitHash,
		sourceBuild.ImageName,
		"",
	)
	if err != nil {
		return err
	}

	return ea.queue.PublishDeployImage(&coreValue.ImageDeployment{
		DeploymentId:          deployment.Id,
		DeploymentName:        normalizedDeploymentName,
		Namespace:             env.Namespace,
		KubeconfigBase64:      kubeconfigBase64,
		ImageRegistryUrl:      ea.cfg.ImageRegistryUrl,
		ImageRegistryUsername: ea.cfg.ImageRegistryUsername,
		ImageRegistryPassword: ea.cfg.ImageRegistryPassword,
		ImageName:             *sourceBuild.ImageName,
		ImageTag:              *sourceBuild.CommitHash,
		Port:                  deploymentPort,
		EnvVars:               coreEnvs,
	})
}

func (da *DeploymentApplication) HandleEnvironmentNotification(notification *coreValue.EnvironmentNotification) {
	environmentId, err := da.GetDeploymentEnvironmentId(notification.DeploymentId)
	if err != nil {
		log.Printf("failed to get environment id for deployment %d: %v", notification.DeploymentId, err)
		return
	}
	if environmentId == nil {
		log.Printf("deployment %d has no environment id", notification.DeploymentId)
		return
	}

	da.notificationHub.Broadcast(notification.CorrelationId, *environmentId, notification)
}
