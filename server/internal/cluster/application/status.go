package application

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

type StatusApplication struct {
	health port.Health
	pubsub port.Pubsub
}

func NewStatusApplication(
	health port.Health,
	pubsub port.Pubsub,
) *StatusApplication {
	return &StatusApplication{
		health: health,
		pubsub: pubsub,
	}
}

func (sa *StatusApplication) HandleRequestDeploymentStatus(d *value.DatabaseDeployment) {
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

	releaseName := d.DeploymentName
	health, err := sa.health.CheckPodsHealthy(releaseName, kubeconfigPath)
	if err != nil {
		log.Printf("failed to check pods health: %v\n", err)
	}

	err = sa.pubsub.PublishDeploymentStatusResponse(&value.HealthStatus{
		DeploymentId: d.DeploymentId,
		Health:       value.Health(health.Health),
		Status:       health.Status,
	})
	if err != nil {
		log.Printf("failed to publish deployment status: %v\n", err)
	}
}
