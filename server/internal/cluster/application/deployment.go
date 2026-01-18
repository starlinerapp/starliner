package application

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"starliner.app/internal/cluster/domain/port"
	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/repository/interface"
	"starliner.app/internal/core/domain/value"
	"strconv"
)

type DeploymentApplication struct {
	clusterRepository _interface.ClusterRepository
	deploy            port.Deploy
	crypto            corePort.Crypto
}

func NewDeploymentApplication(
	clusterRepository _interface.ClusterRepository,
	deploy port.Deploy,
	crypto corePort.Crypto,
) *DeploymentApplication {
	return &DeploymentApplication{
		clusterRepository: clusterRepository,
		deploy:            deploy,
		crypto:            crypto,
	}
}

func (da *DeploymentApplication) HandleDeployDatabase(d *value.Deployment) {
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

	err = da.deploy.DeployPostgres(strconv.FormatInt(d.DeploymentId, 10), kubeconfigPath)
	if err != nil {
		fmt.Printf("failed to install helm chart: %v\n", err)
		return
	}
	log.Println("successfully deployed database to cluster")
}
