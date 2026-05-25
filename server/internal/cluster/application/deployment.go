package application

import (
	"fmt"
	"log"

	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

type DeploymentApplication struct {
	deploy   port.Deploy
	queue    port.Queue
	notifier *Notifier
}

func NewDeploymentApplication(deploy port.Deploy, queue port.Queue) *DeploymentApplication {
	notifier := NewNotifier(queue)
	return &DeploymentApplication{
		deploy:   deploy,
		queue:    queue,
		notifier: notifier,
	}
}

func (da *DeploymentApplication) HandleDeleteDeployment(d *value.Deployment) {
	if d.CorrelationId == nil {
		log.Printf("missing correlation id for deployment %d\n", d.DeploymentId)
	}

	releaseName := d.DeploymentName
	err := da.deploy.DeleteDeployment(d.Namespace, releaseName, d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to delete helm chart: %v\n", err)
		da.notifier.publishNotification(d.DeploymentId, *d.CorrelationId, "failed", fmt.Sprintf("Failed to delete service: %s", d.DeploymentName))
		return
	}
	log.Println("successfully deleted deployment")
	da.notifier.publishNotification(d.DeploymentId, *d.CorrelationId, "success", fmt.Sprintf("Deleted deployment: %s", d.DeploymentName))

	err = da.queue.PublishDeploymentDeleted(&value.DeploymentDeleted{
		DeploymentId: d.DeploymentId,
	})
	if err != nil {
		log.Printf("failed to publish event: %v\n", err)
	}
}
