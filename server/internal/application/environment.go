package application

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"starliner.app/internal/domain/port"
	"starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/domain/service"
	"starliner.app/internal/domain/value"
	"strings"
)

type EnvironmentApplication struct {
	organizationService   *service.OrganizationService
	environmentService    *service.EnvironmentService
	environmentRepository interfaces.EnvironmentRepository
	clusterRepository     interfaces.ClusterRepository
	queue                 port.Queue
	deploy                port.Deploy
	crypto                port.Crypto
}

func NewEnvironmentApplication(
	environmentRepository interfaces.EnvironmentRepository,
	clusterRepository interfaces.ClusterRepository,
	organizationService *service.OrganizationService,
	environmentService *service.EnvironmentService,
	queue port.Queue,
	deploy port.Deploy,
	crypto port.Crypto,
) *EnvironmentApplication {
	return &EnvironmentApplication{
		environmentRepository: environmentRepository,
		clusterRepository:     clusterRepository,
		organizationService:   organizationService,
		environmentService:    environmentService,
		queue:                 queue,
		deploy:                deploy,
		crypto:                crypto,
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

	trimmed := strings.TrimSpace(name)
	environmentSlug := strings.ReplaceAll(strings.ToLower(trimmed), " ", "-")

	_, err = ea.environmentRepository.CreateEnvironment(ctx, name, environmentSlug, projectId)
	if err != nil {
		return err
	}
	return nil
}

func (ea *EnvironmentApplication) DeployDatabase(
	ctx context.Context,
	userId int64,
	environmentId int64,
	database value.Database,
) error {
	err := ea.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return err
	}

	cluster, err := ea.environmentRepository.GetEnvironmentCluster(ctx, environmentId)
	if err != nil {
		return err
	}

	err = ea.queue.PublishDeployDatabase(&value.DeploymentMessage{
		ClusterId: cluster.Id,
		Database:  database,
	})
	if err != nil {
		log.Printf("error publishing: %v", err)
	}
	return nil
}

func (ea *EnvironmentApplication) HandleDeployDatabase(d *value.DeploymentMessage) {
	ctx := context.Background()
	cluster, err := ea.clusterRepository.GetCluster(ctx, d.ClusterId)
	if err != nil {
		fmt.Printf("failed to get cluster from database: %v\n", err)
	}

	kubeconfigBase64, err := ea.crypto.Decrypt(*cluster.Kubeconfig)
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

	err = ea.deploy.DeployPostgres(kubeconfigPath)
	if err != nil {
		fmt.Printf("failed to install helm chart: %v\n", err)
		return
	}
	fmt.Printf("Successfully deployed database to cluster.")
}
