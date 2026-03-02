package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	corePort "starliner.app/internal/core/domain/port"
	coreValue "starliner.app/internal/core/domain/value"
	"strconv"
)

type DeploymentApplication struct {
	environmentService    *service.EnvironmentService
	deploymentService     *service.DeploymentService
	environmentRepository interfaces.EnvironmentRepository
	deploymentRepository  interfaces.DeploymentRepository
	queue                 port.Queue
	pubsub                port.Pubsub
	crypto                corePort.Crypto
}

func NewDeploymentApplication(
	environmentService *service.EnvironmentService,
	deploymentService *service.DeploymentService,
	environmentRepository interfaces.EnvironmentRepository,
	deploymentRepository interfaces.DeploymentRepository,
	queue port.Queue,
	pubsub port.Pubsub,
	crypto corePort.Crypto,
) *DeploymentApplication {
	return &DeploymentApplication{
		environmentService:    environmentService,
		deploymentService:     deploymentService,
		environmentRepository: environmentRepository,
		deploymentRepository:  deploymentRepository,
		queue:                 queue,
		pubsub:                pubsub,
		crypto:                crypto,
	}
}

func (da *DeploymentApplication) DeployImage(
	ctx context.Context,
	userId int64,
	environmentId int64,
	serviceName string,
	imageName string,
	tag string,
	port int,
	envs []*value.EnvVar,
) error {
	err := da.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return err
	}

	cluster, err := da.environmentRepository.GetEnvironmentCluster(ctx, environmentId)
	if err != nil {
		return err
	}

	// TODO: status shouldn't be hardcoded
	deployment, err := da.deploymentRepository.CreateImageDeployment(
		ctx,
		serviceName,
		imageName,
		tag,
		strconv.Itoa(port),
		"unhealthy",
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

	err = da.queue.PublishDeployImage(&coreValue.ImageDeployment{
		DeploymentId:     deployment.Id,
		DeploymentName:   serviceName,
		KubeconfigBase64: kubeconfigBase64,
		ImageRepository:  imageName,
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
	database value.Database,
) error {
	err := da.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return err
	}

	cluster, err := da.environmentRepository.GetEnvironmentCluster(ctx, environmentId)
	if err != nil {
		return err
	}

	// TODO: Replace with real values
	deployment, err := da.deploymentRepository.CreateDatabaseDeployment(
		ctx,
		fmt.Sprintf("%s-%s", string(database), uuid.New().String()[:8]),
		"5432",
		"unhealthy",
		"postgres",
		"postgres",
		environmentId,
	)
	if err != nil {
		return err
	}

	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}

	err = da.queue.PublishDeployDatabase(&coreValue.Deployment{
		DeploymentId:     deployment.Id,
		DeploymentName:   deployment.Name,
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

	cluster, err := da.environmentRepository.GetEnvironmentCluster(ctx, environmentId)
	if err != nil {
		return err
	}

	ingressDeployment, err := da.deploymentRepository.CreateIngressDeployment(
		ctx,
		fmt.Sprintf("ingress-%s", uuid.New().String()[:8]),
		"80",
		"unhealthy",
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

			ch.Paths = append(ch.Paths, coreValue.IngressPath{
				Path:        p.Path,
				PathType:    coreValue.PathType(p.PathType),
				ServiceName: p.ServiceName,
				ServicePort: targetPort,
			})
		}

		coreHosts = append(coreHosts, ch)
	}

	err = da.queue.PublishDeployIngress(&coreValue.IngressDeployment{
		IngressHosts:     coreHosts,
		DeploymentId:     ingressDeployment.Id,
		DeploymentName:   ingressDeployment.Name,
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

	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		return err
	}
	err = da.queue.PublishDeleteDeployment(&coreValue.Deployment{
		DeploymentId:     deployment.Id,
		DeploymentName:   deployment.Name,
		KubeconfigBase64: kubeconfigBase64,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
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

			err = da.pubsub.PublishDeploymentStatusRequest(&coreValue.Deployment{
				DeploymentId:     d.Deployment.Id,
				DeploymentName:   d.Deployment.Name,
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
