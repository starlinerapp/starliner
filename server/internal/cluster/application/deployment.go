package application

import (
	"log"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

type DeploymentApplication struct {
	deploy port.Deploy
	queue  port.Queue
}

func NewDeploymentApplication(deploy port.Deploy, queue port.Queue) *DeploymentApplication {
	return &DeploymentApplication{
		deploy: deploy,
		queue:  queue,
	}
}

func (da *DeploymentApplication) HandleDeleteDeployment(d *value.Deployment) {
	releaseName := d.DeploymentName
	err := da.deploy.DeletePostgres(releaseName, d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to delete helm chart: %v\n", err)
	}
	log.Println("successfully deleted database from cluster")

	err = da.queue.PublishDeploymentDeleted(&value.DeploymentDeleted{
		DeploymentId: d.DeploymentId,
	})
	if err != nil {
		log.Printf("failed to publish event: %v\n", err)
	}
}
