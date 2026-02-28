package application

import (
	"log"
	"starliner.app/internal/cluster/domain/port"
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

func (da *DatabaseApplication) HandleDeployDatabase(d *value.Deployment) {
	// TODO: Check if Cluster CRD is already installed, install otherwise
	err := da.deploy.DeployCloudNativePg("cloudnative-pg", d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to deploy cloudnative-pg: %v\n", err)
	}

	releaseName := d.DeploymentName
	err = da.deploy.DeployPostgres(releaseName, d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to deploy database: %v\n", err)
	}
}
