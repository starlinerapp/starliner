package application

import (
	"context"
	"log"
	interfaces2 "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

type DeploymentApplication struct {
	environmentService    *service.EnvironmentService
	environmentRepository interfaces2.EnvironmentRepository
	deploymentRepository  interfaces2.DeploymentRepository
	queue                 port.Queue
}

func NewDeploymentApplication(
	environmentService *service.EnvironmentService,
	environmentRepository interfaces2.EnvironmentRepository,
	deploymentRepository interfaces2.DeploymentRepository,
	queue port.Queue,
) *DeploymentApplication {
	return &DeploymentApplication{
		environmentService:    environmentService,
		environmentRepository: environmentRepository,
		deploymentRepository:  deploymentRepository,
		queue:                 queue,
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

	err = da.queue.PublishDeployDatabase(&value.Deployment{
		DeploymentId: deployment.Id,
		ClusterId:    cluster.Id,
		Database:     database,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}
