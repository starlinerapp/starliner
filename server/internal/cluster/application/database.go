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
	correlationId := ""
	if d.CorrelationId != nil {
		correlationId = *d.CorrelationId
	} else {
		log.Printf("missing correlation id for DB deployment %d\n", d.DeploymentId)
	}

	releaseName := d.DeploymentName
	err := da.deploy.DeployPostgres(d.Namespace, releaseName, d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to deploy database: %v\n", err)
		if pubErr := da.queue.PublishDatabaseDeployedFailure(&value.DatabaseDeployedFailure{
			CorrelationId:  correlationId,
			DeploymentId:   d.DeploymentId,
			DeploymentName: d.DeploymentName,
		}); pubErr != nil {
			log.Printf("failed to publish database deployed failure: %v\n", pubErr)
		}
		return
	}

	if pubErr := da.queue.PublishDatabaseDeployedSuccess(&value.DatabaseDeployedSuccess{
		CorrelationId:  correlationId,
		DeploymentId:   d.DeploymentId,
		DeploymentName: d.DeploymentName,
		DbName:         "postgres",
		Username:       "postgres",
		Password:       "postgres",
	}); pubErr != nil {
		log.Printf("failed to publish database deployed success: %v\n", pubErr)
	}
}
