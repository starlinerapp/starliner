package application

import (
	"log"
	"starliner.app/internal/cluster/domain/port"
	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

type DatabaseApplication struct {
	deploy port.Deploy
	secret port.Secret
	health port.Health
	queue  port.Queue
	pubsub port.Pubsub
	crypto corePort.Crypto
}

func NewDatabaseApplication(
	deploy port.Deploy,
	health port.Health,
	secret port.Secret,
	queue port.Queue,
	pubsub port.Pubsub,
	crypto corePort.Crypto,
) *DatabaseApplication {
	return &DatabaseApplication{
		deploy: deploy,
		health: health,
		secret: secret,
		queue:  queue,
		pubsub: pubsub,
		crypto: crypto,
	}
}

func (da *DatabaseApplication) HandleDeployDatabase(d *value.Deployment) {
	err := da.deploy.DeployCloudNativePg(d.Namespace, "cloudnative-pg", d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to deploy cloudnative-pg: %v\n", err)
	}

	releaseName := d.DeploymentName
	err = da.deploy.DeployPostgres(d.Namespace, releaseName, d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to deploy database: %v\n", err)
	}

	credentials, err := da.secret.GetDatabaseCredentials(releaseName, d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to get database credentials: %v\n", err)
	}

	err = da.queue.PublishDatabaseDeployed(&value.DatabaseDeployment{
		DeploymentId: d.DeploymentId,
		DbName:       credentials.DatabaseName,
		Username:     credentials.Username,
		Password:     credentials.Password,
	})
	if err != nil {
		log.Printf("failed to publish event: %v\n", err)
	}
}
