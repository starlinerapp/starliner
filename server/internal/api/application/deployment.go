package application

import (
	"context"
	"log"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	corePort "starliner.app/internal/core/domain/port"
	coreValue "starliner.app/internal/core/domain/value"
)

type DeploymentApplication struct {
	environmentService    *service.EnvironmentService
	deploymentService     *service.DeploymentService
	environmentRepository interfaces.EnvironmentRepository
	deploymentRepository  interfaces.DeploymentRepository
	queue                 port.Queue
	crypto                corePort.Crypto
}

func NewDeploymentApplication(
	environmentService *service.EnvironmentService,
	deploymentService *service.DeploymentService,
	environmentRepository interfaces.EnvironmentRepository,
	deploymentRepository interfaces.DeploymentRepository,
	queue port.Queue,
	crypto corePort.Crypto,
) *DeploymentApplication {
	return &DeploymentApplication{
		environmentService:    environmentService,
		deploymentService:     deploymentService,
		environmentRepository: environmentRepository,
		deploymentRepository:  deploymentRepository,
		queue:                 queue,
		crypto:                crypto,
	}
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
		string(database),
		"5432",
		"postgres",
		"test",
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
		KubeconfigBase64: kubeconfigBase64,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (da *DeploymentApplication) DeleteDatabase(ctx context.Context, deploymentId int64, userId int64) error {
	err := da.deploymentService.ValidateUserPermission(ctx, userId, deploymentId)
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
	err = da.queue.PublishDeleteDatabase(&coreValue.Deployment{
		DeploymentId:     deploymentId,
		KubeconfigBase64: kubeconfigBase64,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}

func (da *DeploymentApplication) HandleDatabaseDeleted(c *coreValue.DeploymentDeleted) {
	ctx := context.Background()
	err := da.deploymentRepository.DeleteDeployment(ctx, c.DeploymentId)
	if err != nil {
		log.Printf("failed to delete deployment from database: %v\n", err)
	}
}
