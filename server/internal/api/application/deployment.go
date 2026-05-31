package application

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
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

type DeploymentApplication struct {
	config                *conf.Config
	environmentService    *service.EnvironmentService
	deploymentService     *service.DeploymentService
	parserService         *service.ParserService
	resolverService       *service.ResolverService
	normalizerService     *coreService.NormalizerService
	environmentRepository interfaces.EnvironmentRepository
	deploymentRepository  interfaces.DeploymentRepository
	buildRepository       interfaces.BuildRepository
	githubAppRepository   interfaces.GithubAppRepository
	gitHub                port.GitHub
	grpcClusterClient     port.ClusterClient
	queue                 port.Queue
	pubsub                port.Pubsub
	crypto                corePort.Crypto
}

func NewDeploymentApplication(
	config *conf.Config,
	environmentService *service.EnvironmentService,
	deploymentService *service.DeploymentService,
	parserService *service.ParserService,
	resolverService *service.ResolverService,
	normalizerService *coreService.NormalizerService,
	environmentRepository interfaces.EnvironmentRepository,
	deploymentRepository interfaces.DeploymentRepository,
	buildRepository interfaces.BuildRepository,
	githubAppRepository interfaces.GithubAppRepository,
	gitHub port.GitHub,
	grpcClusterClient port.ClusterClient,
	queue port.Queue,
	pubsub port.Pubsub,
	crypto corePort.Crypto,
) *DeploymentApplication {
	return &DeploymentApplication{
		config:                config,
		environmentService:    environmentService,
		deploymentService:     deploymentService,
		parserService:         parserService,
		resolverService:       resolverService,
		normalizerService:     normalizerService,
		environmentRepository: environmentRepository,
		deploymentRepository:  deploymentRepository,
		buildRepository:       buildRepository,
		githubAppRepository:   githubAppRepository,
		gitHub:                gitHub,
		grpcClusterClient:     grpcClusterClient,
		queue:                 queue,
		pubsub:                pubsub,
		crypto:                crypto,
	}
}

func (da *DeploymentApplication) DeployFromGit(
	ctx context.Context,
	userId int64,
	environmentId int64,
	serviceName string,
	port int,
	gitUrl string,
	projectRepositoryPath string,
	dockerfilePath string,
	envs []*value.EnvVar,
	args []*value.Arg,
) error {
	err := da.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return err
	}

	found, err := da.deploymentRepository.GetEnvironmentDeploymentByName(ctx, environmentId, serviceName)

	if err != nil {
		return err
	}
	if found != nil {
		return fmt.Errorf("%w: %s", value.ErrDeploymentNameAlreadyExists, serviceName)
	}

	env, err := da.environmentRepository.GetEnvironmentById(ctx, environmentId)
	if err != nil {
		return err
	}

	d, err := da.deploymentRepository.CreateGitDeployment(
		ctx,
		environmentId,
		serviceName,
		strconv.Itoa(port),
		gitUrl,
		projectRepositoryPath,
		dockerfilePath,
		envs,
		args,
	)
	if err != nil {
		return err
	}

	b, err := da.buildRepository.CreateBuild(ctx, d.Id, "manual")
	if err != nil {
		return err
	}

	normalizedServiceName, err := da.normalizerService.FormatToDNS1123(serviceName)
	if err != nil {
		return err
	}

	ghApp, err := da.githubAppRepository.GetEnvironmentGithubApp(ctx, environmentId)
	if err != nil {
		return err
	}

	accessToken, err := da.gitHub.GetInstallationToken(ctx, ghApp.InstallationID)
	if err != nil {
		return err
	}

	coreArgs := make([]*coreValue.Arg, len(args))
	for i, a := range args {
		coreArgs[i] = &coreValue.Arg{
			Name:  a.Name,
			Value: a.Value,
		}
	}

	return da.queue.PublishBuildTriggered(&coreValue.TriggerBuild{
		BuildId:        b.Id,
		DeploymentId:   d.Id,
		ImageName:      fmt.Sprintf("%s/%s", env.Namespace, normalizedServiceName),
		GitUrl:         gitUrl,
		BranchName:     env.ConnectedBranch,
		AccessToken:    accessToken,
		RootDirectory:  projectRepositoryPath,
		DockerfilePath: dockerfilePath,
		Args:           coreArgs,
	})
}

func (da *DeploymentApplication) UpdateDeployFromGit(
	ctx context.Context,
	userId int64,
	environmentId int64,
	deploymentId int64,
	port int,
	projectRepositoryPath string,
	dockerfilePath string,
	envs []*value.EnvVar,
	args []*value.Arg,
) (int64, error) {
	err := da.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return 0, err
	}

	existing, err := da.deploymentRepository.GetUserGitDeploymentById(ctx, userId, deploymentId)
	if err != nil {
		return 0, err
	}
	if existing.EnvironmentId == nil || *existing.EnvironmentId != environmentId {
		return 0, fmt.Errorf("git deployment not found")
	}

	env, err := da.environmentRepository.GetEnvironmentById(ctx, environmentId)
	if err != nil {
		return 0, err
	}

	d, err := da.redeployGitDeployment(
		ctx,
		existing,
		strconv.Itoa(port),
		projectRepositoryPath,
		dockerfilePath,
		envs,
		args,
	)
	if err != nil {
		return 0, err
	}

	b, err := da.buildRepository.CreateBuild(ctx, d.Id, "manual")
	if err != nil {
		return 0, err
	}

	normalizedServiceName, err := da.normalizerService.FormatToDNS1123(d.Name)
	if err != nil {
		return 0, err
	}

	ghApp, err := da.githubAppRepository.GetEnvironmentGithubApp(ctx, environmentId)
	if err != nil {
		return 0, err
	}
	accessToken, err := da.gitHub.GetInstallationToken(ctx, ghApp.InstallationID)
	if err != nil {
		return 0, err
	}

	coreArgs := make([]*coreValue.Arg, len(args))
	for i, a := range args {
		coreArgs[i] = &coreValue.Arg{
			Name:  a.Name,
			Value: a.Value,
		}
	}

	err = da.queue.PublishBuildTriggered(&coreValue.TriggerBuild{
		BuildId:        b.Id,
		DeploymentId:   d.Id,
		ImageName:      fmt.Sprintf("%s/%s", env.Namespace, normalizedServiceName),
		AccessToken:    accessToken,
		GitUrl:         d.GitUrl,
		BranchName:     env.ConnectedBranch,
		RootDirectory:  projectRepositoryPath,
		DockerfilePath: dockerfilePath,
		Args:           coreArgs,
	})
	if err != nil {
		return 0, err
	}

	return d.Id, nil
}

func (da *DeploymentApplication) redeployGitDeployment(
	ctx context.Context,
	existing *entity.GitDeployment,
	port string,
	projectRepositoryPath string,
	dockerfilePath string,
	envs []*value.EnvVar,
	args []*value.Arg,
) (*entity.GitDeployment, error) {
	if existing.EnvironmentId == nil {
		return nil, fmt.Errorf("deployment %d has nil environment id", existing.Id)
	}

	if err := da.deploymentRepository.SoftDeleteDeployment(ctx, existing.Id); err != nil {
		return nil, err
	}

	newDeployment, err := da.deploymentRepository.CreateGitDeployment(
		ctx,
		*existing.EnvironmentId,
		existing.Name,
		port,
		existing.GitUrl,
		projectRepositoryPath,
		dockerfilePath,
		envs,
		args,
	)
	if err != nil {
		return nil, err
	}

	if err := da.deploymentRepository.RepointIngressPathsTargetDeployment(ctx, existing.Id, newDeployment.Id); err != nil {
		return nil, err
	}

	return newDeployment, nil
}

func (da *DeploymentApplication) DeployImage(
	ctx context.Context,
	userId int64,
	environmentId int64,
	serviceName string,
	imageName string,
	tag string,
	port int,
	volumeSizeMiB *int32,
	volumeMountPath *string,
	envs []*value.EnvVar,
) error {
	err := da.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return err
	}

	found, err := da.deploymentRepository.GetEnvironmentDeploymentByName(ctx, environmentId, serviceName)

	if err != nil {
		return err
	}
	if found != nil {
		return fmt.Errorf("%w: %s", value.ErrDeploymentNameAlreadyExists, serviceName)
	}

	cluster, err := da.environmentRepository.GetEnvironmentCluster(ctx, environmentId)
	if err != nil {
		return err
	}

	env, err := da.environmentRepository.GetEnvironmentById(ctx, environmentId)
	if err != nil {
		return err
	}

	if (volumeSizeMiB != nil && volumeMountPath == nil) || (volumeSizeMiB == nil && volumeMountPath != nil) {
		return fmt.Errorf("volumeMountPath=%v, volumeSizeMiB=%v: must be both nil or both not nil", volumeMountPath, volumeSizeMiB)
	}

	deployment, err := da.deploymentRepository.CreateImageDeployment(
		ctx,
		serviceName,
		imageName,
		tag,
		strconv.Itoa(port),
		volumeSizeMiB,
		volumeMountPath,
		environmentId,
		envs,
	)
	if err != nil {
		return err
	}

	if cluster.Kubeconfig == nil {
		return fmt.Errorf("cluster kubeconfig is nil")
	}
	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	coreEnvs := make([]*coreValue.EnvVar, 0, len(envs))
	for _, e := range envs {
		res, err := da.parserService.Parse(e.Value)
		if err != nil {
			log.Printf("failed to parse env var: %v\n", err)
			continue
		}

		if deployment.EnvironmentId == nil {
			log.Printf("deployment %d has nil environment id", deployment.Id)
			continue
		}
		resolvedValue, err := da.resolverService.Resolve(ctx, *deployment.EnvironmentId, res)
		if err != nil {
			log.Printf("failed to resolve env var: %v\n", err)
			continue
		}

		coreEnvs = append(coreEnvs, &coreValue.EnvVar{
			Name:  e.Name,
			Value: resolvedValue,
		})
	}

	normalizedServiceName, err := da.normalizerService.FormatToDNS1123(serviceName)
	if err != nil {
		return err
	}

	err = da.queue.PublishDeployImage(&coreValue.ImageDeployment{
		DeploymentId:          deployment.Id,
		DeploymentName:        normalizedServiceName,
		KubeconfigBase64:      kubeconfigBase64,
		Namespace:             env.Namespace,
		ImageRegistryUrl:      da.config.ImageRegistryUrl,
		ImageRegistryUsername: da.config.ImageRegistryUsername,
		ImageRegistryPassword: da.config.ImageRegistryPassword,
		ImageName:             imageName,
		ImageTag:              tag,
		Port:                  port,
		VolumeSizeMiB:         volumeSizeMiB,
		VolumeMountPath:       volumeMountPath,
		EnvVars:               coreEnvs,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (da *DeploymentApplication) UpdateImageDeployment(
	ctx context.Context,
	userId int64,
	deploymentId int64,
	environmentId int64,
	imageName string,
	tag string,
	port int,
	envs []*value.EnvVar) error {
	err := da.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return err
	}

	cluster, err := da.environmentRepository.GetEnvironmentCluster(ctx, environmentId)
	if err != nil {
		return err
	}

	env, err := da.environmentRepository.GetEnvironmentById(ctx, environmentId)
	if err != nil {
		return err
	}

	deployment, err := da.deploymentRepository.UpdateImageDeployment(
		ctx,
		deploymentId,
		imageName,
		tag,
		strconv.Itoa(port),
		envs,
	)
	if err != nil {
		return err
	}

	if deployment.EnvironmentId == nil {
		return fmt.Errorf("deployment %d has nil environment id", deployment.Id)
	}

	if cluster.Kubeconfig == nil {
		return fmt.Errorf("cluster kubeconfig is nil")
	}
	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	coreEnvs := make([]*coreValue.EnvVar, 0, len(envs))
	for _, e := range envs {
		res, err := da.parserService.Parse(e.Value)
		if err != nil {
			log.Printf("failed to parse env var: %v\n", err)
			continue
		}

		resolvedValue, err := da.resolverService.Resolve(ctx, *deployment.EnvironmentId, res)
		if err != nil {
			log.Printf("failed to resolve env var: %v\n", err)
			continue
		}

		coreEnvs = append(coreEnvs, &coreValue.EnvVar{
			Name:  e.Name,
			Value: resolvedValue,
		})
	}

	normalizedServiceName, err := da.normalizerService.FormatToDNS1123(deployment.ServiceName)
	if err != nil {
		return err
	}

	err = da.queue.PublishDeployImage(&coreValue.ImageDeployment{
		DeploymentId:          deployment.Id,
		DeploymentName:        normalizedServiceName,
		Namespace:             env.Namespace,
		KubeconfigBase64:      kubeconfigBase64,
		ImageRegistryUrl:      da.config.ImageRegistryUrl,
		ImageRegistryUsername: da.config.ImageRegistryUsername,
		ImageRegistryPassword: da.config.ImageRegistryPassword,
		ImageName:             imageName,
		ImageTag:              tag,
		Port:                  port,
		VolumeSizeMiB:         deployment.VolumeSizeMiB,
		VolumeMountPath:       deployment.VolumeMountPath,
		EnvVars:               coreEnvs,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (da *DeploymentApplication) DeployDatabase(
	ctx context.Context,
	userId int64,
	environmentId int64,
	serviceName string,
) error {
	err := da.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return err
	}

	found, err := da.deploymentRepository.GetEnvironmentDeploymentByName(ctx, environmentId, serviceName)

	if err != nil {
		return err
	}
	if found != nil {
		return fmt.Errorf("%w: %s", value.ErrDeploymentNameAlreadyExists, serviceName)
	}

	cluster, err := da.environmentRepository.GetEnvironmentCluster(ctx, environmentId)
	if err != nil {
		return err
	}

	env, err := da.environmentRepository.GetEnvironmentById(ctx, environmentId)
	if err != nil {
		return err
	}

	deployment, err := da.deploymentRepository.CreateDatabaseDeployment(
		ctx,
		serviceName,
		"5432",
		environmentId,
	)
	if err != nil {
		return err
	}

	if cluster.Kubeconfig == nil {
		return fmt.Errorf("cluster kubeconfig is nil")
	}
	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	normalizedServiceName, err := da.normalizerService.FormatToDNS1123(serviceName)
	if err != nil {
		return err
	}

	err = da.queue.PublishDeployDatabase(&coreValue.Deployment{
		DeploymentId:     deployment.Id,
		DeploymentName:   normalizedServiceName,
		Namespace:        env.Namespace,
		KubeconfigBase64: kubeconfigBase64,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (da *DeploymentApplication) DeployIngress(ctx context.Context, hosts []*value.IngressHost, userId int64, environmentId int64) error {
	err := da.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return err
	}

	err = da.deploymentService.ValidateIngressHostsAvailable(ctx, hosts)
	if err != nil {
		return err
	}

	cluster, err := da.environmentRepository.GetEnvironmentCluster(ctx, environmentId)
	if err != nil {
		return err
	}

	env, err := da.environmentRepository.GetEnvironmentById(ctx, environmentId)
	if err != nil {
		return err
	}

	ingressDeployment, err := da.deploymentRepository.CreateIngressDeployment(
		ctx,
		fmt.Sprintf("ingress-%s", uuid.New().String()[:8]),
		"80",
		environmentId,
		hosts,
	)
	if err != nil {
		return err
	}

	err = createDeployOnlyBuild(ctx, da.buildRepository, ingressDeployment.Id, value.BuildSourceManual)
	if err != nil {
		return err
	}

	if cluster.Kubeconfig == nil {
		return fmt.Errorf("cluster kubeconfig is nil")
	}
	if cluster.IPv4Address == nil || *cluster.IPv4Address == "" {
		return fmt.Errorf("cluster ipv4 address is not set")
	}
	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	coreHosts := make([]coreValue.IngressHost, 0, len(hosts))
	for _, h := range hosts {
		ch := coreValue.IngressHost{
			Host: h.Host,
		}
		ch.Paths = make([]coreValue.IngressPath, 0, len(h.Paths))

		for _, p := range h.Paths {
			target, err := da.environmentRepository.GetEnvironmentDeploymentByName(ctx, p.ServiceName, environmentId)
			if err != nil {
				return err
			}

			targetPort, err := strconv.Atoi(target.Port)
			if err != nil {
				return err
			}

			normalizedServiceName, err := da.normalizerService.FormatToDNS1123(p.ServiceName)
			if err != nil {
				return err
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

	err = da.queue.PublishDeployIngress(&coreValue.IngressDeployment{
		IngressHosts:     coreHosts,
		DeploymentId:     ingressDeployment.Id,
		DeploymentName:   ingressDeployment.Name,
		Namespace:        env.Namespace,
		KubeconfigBase64: kubeconfigBase64,
		ExpectedIP:       *cluster.IPv4Address,
	})

	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (da *DeploymentApplication) UpdateIngressDeployment(
	ctx context.Context,
	userId int64,
	environmentId int64,
	deploymentId int64,
	hosts []*value.IngressHost,
) (int64, error) {
	err := da.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return 0, err
	}

	existing, err := da.deploymentRepository.GetUserDeployment(ctx, userId, deploymentId)
	if err != nil {
		return 0, err
	}
	if existing.DeletedAt != nil {
		return 0, fmt.Errorf("ingress deployment not found")
	}
	if existing.EnvironmentId == nil || *existing.EnvironmentId != environmentId {
		return 0, fmt.Errorf("ingress deployment not found")
	}

	isIngress, err := da.deploymentRepository.IsIngressDeployment(ctx, deploymentId)
	if err != nil {
		return 0, err
	}
	if !isIngress {
		return 0, fmt.Errorf("ingress deployment not found")
	}

	hostsToValidate := make([]*value.IngressHost, 0, len(hosts))
	for _, h := range hosts {
		if h == nil {
			continue
		}
		existingHost, err := da.deploymentRepository.GetIngressHostByName(ctx, h.Host)
		if err != nil {
			return 0, err
		}
		if existingHost != nil && existingHost.DeploymentId == deploymentId {
			continue
		}
		hostsToValidate = append(hostsToValidate, h)
	}

	err = da.deploymentService.ValidateIngressHostsAvailable(ctx, hostsToValidate)
	if err != nil {
		return 0, err
	}

	ingressDeployment, err := da.redeployIngressDeployment(ctx, deploymentId, environmentId, hosts)
	if err != nil {
		return 0, err
	}

	err = createDeployOnlyBuild(ctx, da.buildRepository, ingressDeployment.Id, value.BuildSourceManual)
	if err != nil {
		return 0, err
	}

	cluster, err := da.environmentRepository.GetEnvironmentCluster(ctx, environmentId)
	if err != nil {
		return 0, err
	}

	env, err := da.environmentRepository.GetEnvironmentById(ctx, environmentId)
	if err != nil {
		return 0, err
	}

	if cluster.Kubeconfig == nil {
		return 0, fmt.Errorf("cluster kubeconfig is nil")
	}
	if cluster.IPv4Address == nil || *cluster.IPv4Address == "" {
		return 0, fmt.Errorf("cluster ipv4 address is not set")
	}
	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return 0, err
	}

	coreHosts := make([]coreValue.IngressHost, 0, len(hosts))
	for _, h := range hosts {
		ch := coreValue.IngressHost{
			Host: h.Host,
		}
		ch.Paths = make([]coreValue.IngressPath, 0, len(h.Paths))

		for _, p := range h.Paths {
			target, err := da.environmentRepository.GetEnvironmentDeploymentByName(ctx, p.ServiceName, environmentId)
			if err != nil {
				return 0, err
			}

			targetPort, err := strconv.Atoi(target.Port)
			if err != nil {
				return 0, err
			}

			normalizedServiceName, err := da.normalizerService.FormatToDNS1123(p.ServiceName)
			if err != nil {
				return 0, err
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

	err = da.queue.PublishDeployIngress(&coreValue.IngressDeployment{
		IngressHosts:     coreHosts,
		DeploymentId:     ingressDeployment.Id,
		DeploymentName:   ingressDeployment.Name,
		Namespace:        env.Namespace,
		KubeconfigBase64: kubeconfigBase64,
		ExpectedIP:       *cluster.IPv4Address,
	})

	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return ingressDeployment.Id, nil
}

func (da *DeploymentApplication) redeployIngressDeployment(
	ctx context.Context,
	deploymentId int64,
	environmentId int64,
	hosts []*value.IngressHost,
) (*entity.IngressDeployment, error) {
	deploymentWithNamespace, err := da.deploymentRepository.GetDeploymentWithNamespace(ctx, deploymentId)
	if err != nil {
		return nil, err
	}

	cluster, err := da.deploymentRepository.GetDeploymentCluster(ctx, deploymentId)
	if err != nil {
		return nil, err
	}

	if cluster.Kubeconfig == nil {
		return nil, fmt.Errorf("cluster kubeconfig is nil")
	}
	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return nil, err
	}

	normalizedDeploymentName, err := da.normalizerService.FormatToDNS1123(deploymentWithNamespace.Name)
	if err != nil {
		return nil, err
	}

	if err := da.deploymentRepository.SoftDeleteDeployment(ctx, deploymentId); err != nil {
		return nil, err
	}

	err = da.queue.PublishDeleteDeployment(&coreValue.Deployment{
		DeploymentId:     deploymentWithNamespace.Id,
		DeploymentName:   normalizedDeploymentName,
		Namespace:        deploymentWithNamespace.Namespace,
		KubeconfigBase64: kubeconfigBase64,
	})
	if err != nil {
		log.Printf("error publishing ingress delete: %v", err)
	}

	return da.deploymentRepository.CreateIngressDeployment(
		ctx,
		fmt.Sprintf("ingress-%s", uuid.New().String()[:8]),
		"80",
		environmentId,
		hosts,
	)
}

func (da *DeploymentApplication) DeleteDeployment(ctx context.Context, deploymentId int64, userId int64) error {
	if err := da.deploymentService.ValidateUserPermission(ctx, userId, deploymentId); err != nil {
		return err
	}
	deployment, err := da.deploymentRepository.GetUserDeployment(ctx, userId, deploymentId)
	if err != nil {
		return err
	}
	if deployment.DeletedAt != nil {
		return fmt.Errorf("deployment already deleted")
	}

	deploymentWithNamespace, err := da.deploymentRepository.GetDeploymentWithNamespace(ctx, deploymentId)
	if err != nil {
		return err
	}

	cluster, err := da.deploymentRepository.GetDeploymentCluster(ctx, deploymentId)
	if err != nil {
		return err
	}

	if cluster.Kubeconfig == nil {
		return fmt.Errorf("cluster kubeconfig is nil")
	}
	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	normalizedDeploymentName, err := da.normalizerService.FormatToDNS1123(deploymentWithNamespace.Name)
	if err != nil {
		return err
	}

	err = da.deploymentRepository.SoftDeleteDeploymentVolume(ctx, deploymentId)
	if err != nil {
		return err
	}

	err = da.queue.PublishDeleteDeployment(&coreValue.Deployment{
		DeploymentId:     deploymentWithNamespace.Id,
		DeploymentName:   normalizedDeploymentName,
		Namespace:        deploymentWithNamespace.Namespace,
		KubeconfigBase64: kubeconfigBase64,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (da *DeploymentApplication) StreamDeploymentLogs(ctx context.Context, userId int64, deploymentId int64, w io.Writer) error {
	err := da.deploymentService.ValidateUserPermission(ctx, userId, deploymentId)
	if err != nil {
		return err
	}

	deployment, err := da.deploymentRepository.GetDeploymentWithNamespace(ctx, deploymentId)
	if err != nil {
		return err
	}

	cluster, err := da.deploymentRepository.GetDeploymentCluster(ctx, deploymentId)
	if err != nil {
		return err
	}

	if cluster.Kubeconfig == nil {
		return fmt.Errorf("cluster kubeconfig is nil")
	}
	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	normalizedDeploymentName, err := da.normalizerService.FormatToDNS1123(deployment.Name)
	if err != nil {
		return err
	}

	logSource, err := da.deploymentLogSource(ctx, deploymentId)
	if err != nil {
		return err
	}

	return da.grpcClusterClient.StreamLogs(
		ctx,
		string(logSource),
		deployment.Namespace,
		normalizedDeploymentName,
		kubeconfigBase64,
		w,
	)
}

func (da *DeploymentApplication) StreamDeploymentStatusLogs(
	ctx context.Context,
	userId int64,
	deploymentId int64,
	w io.Writer,
) error {
	err := da.deploymentService.ValidateUserPermission(ctx, userId, deploymentId)
	if err != nil {
		return err
	}

	deployment, err := da.deploymentRepository.GetUserDeployment(ctx, userId, deploymentId)
	if err != nil {
		return err
	}

	stored, err := da.deploymentRepository.GetDeploymentStatusLogs(ctx, userId, deploymentId)
	if err != nil {
		return err
	}

	isIngress, err := da.deploymentRepository.IsIngressDeployment(ctx, deploymentId)
	if err != nil {
		return err
	}

	hasStoredLogs := stored.Logs != nil && *stored.Logs != ""

	if deployment.DeletedAt != nil {
		if hasStoredLogs {
			_, err := io.WriteString(w, *stored.Logs)
			return err
		}
		_, err := io.WriteString(w, "This deployment has been deleted.\n")
		return err
	}

	if hasStoredLogs {
		_, err := io.WriteString(w, *stored.Logs)
		return err
	}

	deploymentWithNamespace, err := da.deploymentRepository.GetDeploymentWithNamespace(ctx, deploymentId)
	if err != nil {
		return err
	}

	cluster, err := da.deploymentRepository.GetDeploymentCluster(ctx, deploymentId)
	if err != nil {
		return err
	}

	if cluster.Kubeconfig == nil {
		return fmt.Errorf("cluster kubeconfig is nil")
	}
	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	normalizedDeploymentName, err := da.normalizerService.FormatToDNS1123(deploymentWithNamespace.Name)
	if err != nil {
		return err
	}

	if isIngress {
		return da.grpcClusterClient.StreamIngressDeploymentStatusLogs(
			ctx,
			deploymentId,
			deploymentWithNamespace.Namespace,
			normalizedDeploymentName,
			kubeconfigBase64,
			w,
		)
	}

	commitHash := da.resolveDeploymentCommitHash(ctx, deploymentWithNamespace)
	return da.grpcClusterClient.StreamDeploymentStatusLogs(
		ctx,
		deploymentId,
		deploymentWithNamespace.Namespace,
		normalizedDeploymentName,
		kubeconfigBase64,
		commitHash,
		w,
	)
}

func rolloutStatusFromLogs(logs *string) string {
	if logs == nil || *logs == "" {
		return "pending"
	}
	if strings.Contains(*logs, "has failed.") || strings.Contains(*logs, "==> ERROR:") {
		return "failure"
	}
	if strings.Contains(*logs, "is complete.") || strings.Contains(*logs, "Ingress deployed successfully") {
		return "success"
	}
	return "pending"
}

func createDeployOnlyBuild(
	ctx context.Context,
	buildRepository interfaces.BuildRepository,
	deploymentId int64,
	source string,
) error {
	b, err := buildRepository.CreateBuild(ctx, deploymentId, source)
	if err != nil {
		return err
	}

	return buildRepository.UpdateBuild(
		ctx,
		b.Id,
		value.BuildStatusSuccess,
		nil,
		nil,
		"Build skipped",
	)
}

func (da *DeploymentApplication) resolveDeploymentCommitHash(
	ctx context.Context,
	deployment *entity.Deployment,
) string {
	if deployment.EnvironmentId == nil {
		return ""
	}

	build, err := da.buildRepository.GetLatestGitDeploymentBuild(
		ctx,
		*deployment.EnvironmentId,
		deployment.Name,
	)
	if err != nil || build == nil || build.CommitHash == nil {
		return ""
	}

	return *build.CommitHash
}

func (da *DeploymentApplication) deploymentLogSource(ctx context.Context, deploymentId int64) (value.LogSource, error) {
	isIngress, err := da.deploymentRepository.IsIngressDeployment(ctx, deploymentId)
	if err != nil {
		return "", err
	}

	if isIngress {
		return value.LogSourceIngress, nil
	}

	return value.LogSourceWorkload, nil
}

func (da *DeploymentApplication) OpenTTY(
	ctx context.Context,
	userId int64,
	deploymentId int64,
	stdin io.Reader,
	stdout io.Writer,
	sizes <-chan port.TerminalSize,
) error {
	err := da.deploymentService.ValidateUserPermission(ctx, userId, deploymentId)
	if err != nil {
		return err
	}

	deployment, err := da.deploymentRepository.GetDeploymentWithNamespace(ctx, deploymentId)
	if err != nil {
		return err
	}

	cluster, err := da.deploymentRepository.GetDeploymentCluster(ctx, deploymentId)
	if err != nil {
		return err
	}

	if cluster.Kubeconfig == nil {
		return fmt.Errorf("cluster kubeconfig is nil")
	}
	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	normalizedDeploymentName, err := da.normalizerService.FormatToDNS1123(deployment.Name)
	if err != nil {
		return err
	}

	return da.grpcClusterClient.OpenTTY(ctx, deployment.Namespace, normalizedDeploymentName, kubeconfigBase64, stdin, stdout, sizes)
}

func (da *DeploymentApplication) HandleDatabaseDeploymentCreated(c *coreValue.DatabaseDeployment) {
	ctx := context.Background()

	encryptedPassword, err := da.crypto.Encrypt(c.Password)
	if err != nil {
		log.Printf("failed to encrypt database password: %v\n", err)
		return
	}

	err = da.deploymentRepository.UpdateDatabaseDeploymentCredentials(ctx, c.DbName, c.DeploymentId, c.Username, encryptedPassword)
	if err != nil {
		log.Printf("failed to update database deployment credentials: %v\n", err)
	}
}

func (da *DeploymentApplication) HandleDeploymentStatusLogsCompleted(c *coreValue.DeploymentStatusLogsCompleted) {
	ctx := context.Background()
	rolloutStatus := rolloutStatusFromLogs(&c.Logs)
	err := da.deploymentRepository.SetDeploymentStatusLogs(ctx, c.DeploymentId, c.Logs, rolloutStatus)
	if err != nil {
		log.Printf("failed to persist deployment status logs: %v", err)
	}
}

func (da *DeploymentApplication) HandleDeploymentDeleted(c *coreValue.DeploymentDeleted) {
	ctx := context.Background()
	err := da.deploymentRepository.SoftDeleteDeployment(ctx, c.DeploymentId)
	if err != nil {
		log.Printf("failed to soft delete deployment from database: %v\n", err)
	}
}

func (da *DeploymentApplication) RequestDeploymentStatus() error {
	ctx := context.Background()
	deployments, err := da.deploymentRepository.GetAllDeploymentsWithKubeconfig(ctx)
	if err != nil {
		return err
	}

	for _, d := range deployments {
		go func(d *entity.DeploymentWithKubeconfig) {
			if d.Kubeconfig == nil {
				log.Printf("deployment status request skipped: kubeconfig is nil for deployment id=%d\n\n", d.Deployment.Id)
				return
			}
			kubeconfigBase64, err := da.crypto.Decrypt(*d.Kubeconfig)
			if err != nil {
				log.Printf("failed to decrypt kubeconfig: %v\n", err)
				return
			}

			normalizedDeploymentName, err := da.normalizerService.FormatToDNS1123(d.Deployment.Name)
			if err != nil {
				log.Printf("failed to normalize deployment name: %v\n", err)
				return
			}

			deployment := &coreValue.Deployment{
				Namespace:        d.Deployment.Namespace,
				DeploymentId:     d.Deployment.Id,
				DeploymentName:   normalizedDeploymentName,
				KubeconfigBase64: kubeconfigBase64,
				ClusterId:        d.ClusterId,
				OrganizationId:   d.OrganizationId,
			}
			if d.ProvisioningId != nil {
				deployment.ProvisioningId = *d.ProvisioningId
			}

			err = da.pubsub.PublishDeploymentStatusRequest(deployment)
			if err != nil {
				log.Printf("failed to publish: %v\n", err)
			}
		}(d)
	}
	return nil
}

func (da *DeploymentApplication) HandleDeploymentStatusResponse(health *coreValue.HealthStatus) {
	ctx := context.Background()
	err := da.deploymentRepository.UpdateDeploymentStatus(ctx, health.DeploymentId, string(health.Health))
	if err != nil {
		log.Printf("failed to update deployment status: %v\n", err)
	}
}

func (da *DeploymentApplication) HandleBuildCompleted(b *coreValue.BuildCompleted) {
	ctx := context.Background()
	err := da.buildRepository.UpdateBuild(ctx, b.BuildId, value.BuildStatus(b.BuildStatus), b.CommitHash, b.ImageName, b.Logs)
	if err != nil {
		log.Printf("failed to update build status: %v\n", err)
	}
	if b.BuildStatus == coreValue.BuildStatusFailed {
		return
	}

	cluster, err := da.deploymentRepository.GetDeploymentCluster(ctx, b.DeploymentId)
	if err != nil {
		log.Printf("failed to get deployment cluster: %v\n", err)
		return
	}

	deployment, err := da.deploymentRepository.GetDeploymentWithNamespace(ctx, b.DeploymentId)
	if err != nil {
		log.Printf("failed to get deployment: %v\n", err)
		return
	}
	if deployment.EnvironmentId == nil {
		log.Printf("deployment %d has nil environment id\n", b.DeploymentId)
		return
	}

	envs, err := da.deploymentRepository.GetDeploymentEnvs(ctx, b.DeploymentId)
	if err != nil {
		log.Printf("failed to get deployment envs: %v\n", err)
		return
	}
	coreEnvs := make([]*coreValue.EnvVar, 0, len(envs))
	for _, e := range envs {
		res, err := da.parserService.Parse(e.Value)
		if err != nil {
			log.Printf("failed to parse env var: %v\n", err)
			continue
		}

		resolvedValue, err := da.resolverService.Resolve(ctx, *deployment.EnvironmentId, res)
		if err != nil {
			log.Printf("failed to resolve env var: %v\n", err)
			continue
		}

		coreEnvs = append(coreEnvs, &coreValue.EnvVar{
			Name:  e.Name,
			Value: resolvedValue,
		})
	}

	if cluster.Kubeconfig == nil {
		log.Printf("cluster kubeconfig is nil for deployment %d\n", b.DeploymentId)
		return
	}
	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		log.Printf("failed to decrypt kubeconfig: %v\n", err)
		return
	}

	deploymentPort, err := strconv.Atoi(deployment.Port)
	if err != nil {
		log.Printf("failed to parse port: %v\n", err)
		return
	}

	normalizedDeploymentName, err := da.normalizerService.FormatToDNS1123(deployment.Name)
	if err != nil {
		log.Printf("failed to normalize deployment name: %v\n", err)
	}

	err = da.queue.PublishDeployImage(&coreValue.ImageDeployment{
		DeploymentId:          b.DeploymentId,
		DeploymentName:        normalizedDeploymentName,
		Namespace:             deployment.Namespace,
		KubeconfigBase64:      kubeconfigBase64,
		ImageRegistryUrl:      da.config.ImageRegistryUrl,
		ImageRegistryUsername: da.config.ImageRegistryUsername,
		ImageRegistryPassword: da.config.ImageRegistryPassword,
		ImageName:             *b.ImageName,
		ImageTag:              *b.Tag,
		Port:                  deploymentPort,
		EnvVars:               coreEnvs,
	})
	if err != nil {
		log.Printf("failed to publish: %v\n", err)
	}
}
