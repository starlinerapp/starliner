package application

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"starliner.app/internal/domain/port"
	interfaces "starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/domain/service"
	"starliner.app/internal/domain/value"
)

type DeploymentApplication struct {
	environmentService    *service.EnvironmentService
	environmentRepository interfaces.EnvironmentRepository
	clusterRepository     interfaces.ClusterRepository
	deploymentRepository  interfaces.DeploymentRepository
	deploy                port.Deploy
	queue                 port.Queue
	crypto                port.Crypto
}

func NewDeploymentApplication(
	environmentService *service.EnvironmentService,
	environmentRepository interfaces.EnvironmentRepository,
	clusterRepository interfaces.ClusterRepository,
	deploymentRepository interfaces.DeploymentRepository,
	deploy port.Deploy,
	queue port.Queue,
	crypto port.Crypto,
) *DeploymentApplication {
	return &DeploymentApplication{
		environmentService:    environmentService,
		environmentRepository: environmentRepository,
		deploymentRepository:  deploymentRepository,
		clusterRepository:     clusterRepository,
		deploy:                deploy,
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

	err = da.queue.PublishDeployDatabase(&value.DeploymentMessage{
		ClusterId: cluster.Id,
		Database:  database,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	err = da.deploymentRepository.CreateDeployment(ctx, string(database), environmentId)
	if err != nil {
		return err
	}
	return nil
}

func (da *DeploymentApplication) HandleDeployDatabase(d *value.DeploymentMessage) {
	ctx := context.Background()
	cluster, err := da.clusterRepository.GetCluster(ctx, d.ClusterId)
	if err != nil {
		fmt.Printf("failed to get cluster from database: %v\n", err)
	}

	kubeconfigBase64, err := da.crypto.Decrypt(*cluster.Kubeconfig)
	if err != nil {
		fmt.Printf("failed to decrypt kubeconfig: %v\n", err)
		return
	}
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(kubeconfigBase64)
	if err != nil {
		fmt.Printf("failed to decode kubeconfig: %v\n", err)
		return
	}

	tmpDir, err := os.MkdirTemp("", "kubeconfig-*")
	if err != nil {
		fmt.Printf("failed to create temp directory: %v\n", err)
		return
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("failed to remove temp directory: %v\n", err)
		}
	}()

	kubeconfigPath := filepath.Join(tmpDir, "kubeconfig")
	err = os.WriteFile(kubeconfigPath, kubeconfigBytes, 0600)
	if err != nil {
		fmt.Printf("failed to write kubeconfig: %v\n", err)
		return
	}

	err = da.deploy.DeployPostgres(kubeconfigPath)
	if err != nil {
		fmt.Printf("failed to install helm chart: %v\n", err)
		return
	}
	fmt.Printf("Successfully deployed database to cluster.")
}
