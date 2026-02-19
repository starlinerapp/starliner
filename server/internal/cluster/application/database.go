package application

import (
	"fmt"
	"log"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/cluster/utils"
	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

type DatabaseApplication struct {
	deploy port.Deploy
	health port.Health
	queue  port.Queue
	pubsub port.Pubsub
	crypto corePort.Crypto
}

func NewDatabaseApplication(
	deploy port.Deploy,
	health port.Health,
	queue port.Queue,
	pubsub port.Pubsub,
	crypto corePort.Crypto,
) *DatabaseApplication {
	return &DatabaseApplication{
		deploy: deploy,
		health: health,
		queue:  queue,
		pubsub: pubsub,
		crypto: crypto,
	}
}

func (da *DatabaseApplication) HandleDeployDatabase(d *value.DatabaseDeployment) {
	err := utils.WithTempKubeConfig(d.KubeconfigBase64, func(kubeconfigPath string) error {
		// TODO: Check if Cluster CRD is already installed, install otherwise
		err := da.deploy.DeployCloudNativePg("cloudnative-pg", kubeconfigPath)
		if err != nil {
			return fmt.Errorf("failed to deploy cloudnative-pg: %v\n", err)
		}

		releaseName := d.DeploymentName
		err = da.deploy.DeployPostgres(releaseName, kubeconfigPath)
		if err != nil {
			return fmt.Errorf("failed to deploy postgres: %v\n", err)
		}

		return nil
	})
	if err != nil {
		log.Printf("failed to deploy database: %v\n", err)
	}
}

func (da *DatabaseApplication) HandleDeleteDatabase(d *value.DatabaseDeployment) {
	err := utils.WithTempKubeConfig(d.KubeconfigBase64, func(kubeconfigPath string) error {
		releaseName := d.DeploymentName
		err := da.deploy.DeletePostgres(releaseName, kubeconfigPath)
		if err != nil {
			return fmt.Errorf("failed to delete helm chart: %v\n", err)
		}
		log.Println("successfully deleted database from cluster")

		err = da.queue.PublishDatabaseDeleted(&value.DeploymentDeleted{
			DeploymentId: d.DeploymentId,
		})
		if err != nil {
			fmt.Printf("failed to publish event: %v\n", err)
		}
		return nil
	})
	if err != nil {
		log.Printf("failed to delete database: %v\n", err)
	}
}
