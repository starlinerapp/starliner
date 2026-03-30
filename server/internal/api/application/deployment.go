package application

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/google/uuid"
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
	environmentService    *service.EnvironmentService
	deploymentService     *service.DeploymentService
	normalizerService     *coreService.NormalizerService
	environmentRepository interfaces.EnvironmentRepository
	deploymentRepository  interfaces.DeploymentRepository
	buildRepository       interfaces.BuildRepository
	grpcClient            port.GrpcClient
	queue                 port.Queue
	pubsub                port.Pubsub
	crypto                corePort.Crypto
}

func NewDeploymentApplication(
	environmentService *service.EnvironmentService,
	deploymentService *service.DeploymentService,
	normalizerService *coreService.NormalizerService,
	environmentRepository interfaces.EnvironmentRepository,
	deploymentRepository interfaces.DeploymentRepository,
	buildRepository interfaces.BuildRepository,
	grpcClient port.GrpcClient,
	queue port.Queue,
	pubsub port.Pubsub,
	crypto corePort.Crypto,
) *DeploymentApplication {
	return &DeploymentApplication{
		environmentService:    environmentService,
		deploymentService:     deploymentService,
		normalizerService:     normalizerService,
		environmentRepository: environmentRepository,
		deploymentRepository:  deploymentRepository,
		buildRepository:       buildRepository,
		grpcClient:            grpcClient,
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
	)
	if err != nil {
		return err
	}

	b, err := da.buildRepository.CreateBuild(ctx, d.Id)
	if err != nil {
		return err
	}

	normalizedServiceName, err := da.normalizerService.FormatToDNS1123(serviceName)
	if err != nil {
		return err
	}

	return da.queue.PublishBuildTriggered(&coreValue.TriggerBuild{
		BuildId:        b.Id,
		DeploymentId:   d.Id,
		ImageName:      fmt.Sprintf("%s/%s", env.Namespace, normalizedServiceName),
		GitUrl:         gitUrl,
		RootDirectory:  projectRepositoryPath,
		DockerfilePath: dockerfilePath,
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
) error {
	err := da.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return err
	}

	env, err := da.environmentRepository.GetEnvironmentById(ctx, environmentId)
	if err != nil {
		return err
	}

	d, err := da.deploymentRepository.UpdateGitDeployment(
		ctx,
		deploymentId,
		strconv.Itoa(port),
		projectRepositoryPath,
		dockerfilePath,
		envs,
	)
	if err != nil {
		return err
	}

	b, err := da.buildRepository.CreateBuild(ctx, d.Id)
	if err != nil {
		return err
	}

	normalizedServiceName, err := da.normalizerService.FormatToDNS1123(d.Name)
	if err != nil {
		return err
	}

	return da.queue.PublishBuildTriggered(&coreValue.TriggerBuild{
		BuildId:        b.Id,
		DeploymentId:   d.Id,
		ImageName:      fmt.Sprintf("%s/%s", env.Namespace, normalizedServiceName),
		GitUrl:         d.GitUrl,
		RootDirectory:  projectRepositoryPath,
		DockerfilePath: dockerfilePath,
	})
}

func (da *DeploymentApplication) DeployImage(
	ctx context.Context,
	userId int64,
	environmentId int64,
	serviceName string,
	imageName string,
	tag string,
	port int,
	volumeSizeMB *int32,
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

	deployment, err := da.deploymentRepository.CreateImageDeployment(
		ctx,
		serviceName,
		imageName,
		tag,
		strconv.Itoa(port),
		volumeSizeMB,
		environmentId,
		envs,
	)
	if err != nil {
		return err
	}

	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	coreEnvs := make([]*coreValue.EnvVar, 0, len(envs))
	for _, e := range envs {
		coreEnvs = append(coreEnvs, &coreValue.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		})
	}

	normalizedServiceName, err := da.normalizerService.FormatToDNS1123(serviceName)
	if err != nil {
		return err
	}

	err = da.queue.PublishDeployImage(&coreValue.ImageDeployment{
		DeploymentId:     deployment.Id,
		DeploymentName:   normalizedServiceName,
		KubeconfigBase64: kubeconfigBase64,
		Namespace:        env.Namespace,
		ImageName:        imageName,
		ImageTag:         tag,
		Port:             port,
		EnvVars:          coreEnvs,
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

	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	coreEnvs := make([]*coreValue.EnvVar, 0, len(envs))
	for _, e := range envs {
		coreEnvs = append(coreEnvs, &coreValue.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		})
	}

	normalizedServiceName, err := da.normalizerService.FormatToDNS1123(deployment.ServiceName)
	if err != nil {
		return err
	}

	err = da.queue.PublishDeployImage(&coreValue.ImageDeployment{
		DeploymentId:     deployment.Id,
		DeploymentName:   normalizedServiceName,
		Namespace:        env.Namespace,
		KubeconfigBase64: kubeconfigBase64,
		ImageName:        imageName,
		ImageTag:         tag,
		Port:             port,
		EnvVars:          coreEnvs,
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
) error {
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

	ingressDeployment, err := da.deploymentRepository.UpdateIngressDeployment(
		ctx,
		deploymentId,
		"80",
		environmentId,
		hosts,
	)
	if err != nil {
		return err
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
	})

	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (da *DeploymentApplication) DeleteDeployment(ctx context.Context, deploymentId int64, userId int64) error {
	deployment, err := da.deploymentRepository.GetUserDeployment(ctx, userId, deploymentId)
	if err != nil {
		return err
	}

	cluster, err := da.deploymentRepository.GetDeploymentCluster(ctx, deploymentId)
	if err != nil {
		return err
	}

	env, err := da.environmentRepository.GetEnvironmentById(ctx, deployment.EnvironmentId)
	if err != nil {
		return err
	}

	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	normalizedDeploymentName, err := da.normalizerService.FormatToDNS1123(deployment.Name)
	if err != nil {
		return err
	}

	err = da.queue.PublishDeleteDeployment(&coreValue.Deployment{
		DeploymentId:     deployment.Id,
		DeploymentName:   normalizedDeploymentName,
		Namespace:        env.Namespace,
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

	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	normalizedDeploymentName, err := da.normalizerService.FormatToDNS1123(deployment.Name)
	if err != nil {
		return err
	}

	return da.grpcClient.StreamLogs(ctx, deployment.Namespace, normalizedDeploymentName, kubeconfigBase64, w)
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

	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	normalizedDeploymentName, err := da.normalizerService.FormatToDNS1123(deployment.Name)
	if err != nil {
		return err
	}

	return da.grpcClient.OpenTTY(ctx, deployment.Namespace, normalizedDeploymentName, kubeconfigBase64, stdin, stdout, sizes)
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

func (da *DeploymentApplication) HandleDeploymentDeleted(c *coreValue.DeploymentDeleted) {
	ctx := context.Background()
	err := da.deploymentRepository.DeleteDeployment(ctx, c.DeploymentId)
	if err != nil {
		log.Printf("failed to delete deployment from database: %v\n", err)
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
			kubeconfigBase64, err := da.crypto.Decrypt(*d.Kubeconfig)
			if err != nil {
				log.Printf("failed to decrypt kubeconfig: %v\n", err)
			}

			normalizedDeploymentName, err := da.normalizerService.FormatToDNS1123(d.Deployment.Name)
			if err != nil {
				log.Printf("failed to normalize deployment name: %v\n", err)
			}

			err = da.pubsub.PublishDeploymentStatusRequest(&coreValue.Deployment{
				Namespace:        d.Deployment.Namespace,
				DeploymentId:     d.Deployment.Id,
				DeploymentName:   normalizedDeploymentName,
				KubeconfigBase64: kubeconfigBase64,
			})
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
	err := da.buildRepository.UpdateBuild(ctx, b.BuildId, value.BuildStatus(b.BuildStatus), b.Logs)
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

	envs, err := da.deploymentRepository.GetDeploymentEnvs(ctx, b.DeploymentId)
	if err != nil {
		log.Printf("failed to get deployment envs: %v\n", err)
		return
	}
	coreEnvs := make([]*coreValue.EnvVar, 0, len(envs))
	for _, e := range envs {
		coreEnvs = append(coreEnvs, &coreValue.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		})
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
		DeploymentId:     b.DeploymentId,
		DeploymentName:   normalizedDeploymentName,
		Namespace:        deployment.Namespace,
		KubeconfigBase64: kubeconfigBase64,
		ImageName:        fmt.Sprintf("%s/%s", b.ImageRegistryUrl, b.ImageName),
		ImageTag:         b.Tag,
		Port:             deploymentPort,
		EnvVars:          coreEnvs,
	})
	if err != nil {
		log.Printf("failed to publish: %v\n", err)
	}
}
