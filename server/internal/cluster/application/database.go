package application

import (
	"fmt"
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
	notifier *Notifier
}

func NewDatabaseApplication(
	deploy port.Deploy,
	health port.Health,
	queue port.Queue,
	pubsub port.Pubsub,
	crypto corePort.Crypto,
) *DatabaseApplication {

	notifier := NewNotifier(queue)
	return &DatabaseApplication{
		deploy: deploy,
		health: health,
		queue:  queue,
		pubsub: pubsub,
		crypto: crypto,
		notifier: notifier,
	}
}

func (da *DatabaseApplication) HandleDeployDatabase(d *value.Deployment) {
	if d.CorrelationId == nil {
		log.Printf("missing correlation id for DB deployment %d\n", d.DeploymentId)
	}

	err := da.deploy.DeployCloudNativePg(d.Namespace, "cloudnative-pg", d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to deploy cloudnative-pg: %v\n", err)
	}

	releaseName := d.DeploymentName
	err := da.deploy.DeployPostgres(d.Namespace, releaseName, d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to deploy database: %v\n", err)
		da.notifier.publishNotification(d.DeploymentId, *d.CorrelationId, "failed", fmt.Sprintf("Failed to deploy database %s", d.DeploymentName))
		return
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

	da.notifier.publishNotification(d.DeploymentId, *d.CorrelationId, "success", fmt.Sprintf("Database %s deployed successfully", d.DeploymentName))
}
