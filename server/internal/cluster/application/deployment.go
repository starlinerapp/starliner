package application

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"starliner.app/internal/cluster/domain/port"
	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

type DeploymentApplication struct {
	deploy port.Deploy
	queue  corePort.Queue
	crypto corePort.Crypto
}

func NewDeploymentApplication(
	deploy port.Deploy,
	queue corePort.Queue,
	crypto corePort.Crypto,
) *DeploymentApplication {
	return &DeploymentApplication{
		deploy: deploy,
		queue:  queue,
		crypto: crypto,
	}
}

func (da *DeploymentApplication) HandleDeployDatabase(d *value.Deployment) {
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(d.KubeconfigBase64)
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

	releaseName := "postgres"
	err = da.deploy.DeployPostgres(releaseName, kubeconfigPath)
	if err != nil {
		fmt.Printf("failed to install helm chart: %v\n", err)
		return
	}
	log.Println("successfully deployed database to cluster")
}

func (da *DeploymentApplication) HandleDeleteDatabase(d *value.Deployment) {
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(d.KubeconfigBase64)
	if err != nil {
		fmt.Printf("failed to decode kubeconfig: %v\n", err)
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

	releaseName := "postgres"
	err = da.deploy.DeletePostgres(releaseName, kubeconfigPath)
	if err != nil {
		fmt.Printf("failed to delete helm chart: %v\n", err)
		return
	}
	log.Println("successfully deleted database from cluster")

	err = da.queue.PublishDatabaseDeleted(&value.DeploymentDeleted{
		DeploymentId: d.DeploymentId,
	})
	if err != nil {
		fmt.Printf("failed to publish event: %v\n", err)
	}
}
