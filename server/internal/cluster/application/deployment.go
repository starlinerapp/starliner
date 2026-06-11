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
	correlationId := ""
	if d.CorrelationId != nil {
		correlationId = *d.CorrelationId
	} else {
		log.Printf("missing correlation id for deployment %d\n", d.DeploymentId)
	}

	releaseName := d.DeploymentName
	err := da.deploy.DeleteDeployment(d.Namespace, releaseName, d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to delete helm chart: %v\n", err)
		if pubErr := da.queue.PublishDeploymentDeletedFailure(&value.DeploymentDeletedFailure{
			CorrelationId:  correlationId,
			DeploymentId:   d.DeploymentId,
			DeploymentName: d.DeploymentName,
		}); pubErr != nil {
			log.Printf("failed to publish deployment deleted failure: %v\n", pubErr)
		}
		return
	}
	log.Println("successfully deleted deployment")

	if pubErr := da.queue.PublishDeploymentDeletedSuccess(&value.DeploymentDeletedSuccess{
		CorrelationId:  correlationId,
		DeploymentId:   d.DeploymentId,
		DeploymentName: d.DeploymentName,
	}); pubErr != nil {
		log.Printf("failed to publish deployment deleted success: %v\n", pubErr)
	}
}
