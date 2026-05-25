package application

import (
	"log"

	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

type Notifier struct {
	queue port.Queue
}

func NewNotifier(queue port.Queue) *Notifier {
	return &Notifier{queue: queue}
}

func (n *Notifier) publishNotification(deploymentId int64, correlationId string, status string, message string) {
	err := n.queue.PublishDeploymentNotification(&value.EnvironmentNotification{
		DeploymentId:  deploymentId,
		CorrelationId: correlationId,
		Status:        status,
		Message:       message,
	})
	if err != nil {
		log.Printf("failed to publish deployment notification: %v\n", err)
	}
}
