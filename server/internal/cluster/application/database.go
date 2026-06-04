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
	releaseName := d.DeploymentName
	err := da.deploy.DeployPostgres(d.Namespace, releaseName, d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to deploy database: %v\n", err)
	}

	err = da.queue.PublishDatabaseDeployed(&value.DatabaseDeployment{
		DeploymentId: d.DeploymentId,
		DbName:       "postgres",
		Username:     "postgres",
		Password:     "postgres",
	})
	if err != nil {
		log.Printf("failed to publish event: %v\n", err)
	}
}
