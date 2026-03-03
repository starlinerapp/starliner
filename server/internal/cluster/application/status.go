package application

import (
	"log"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

type StatusApplication struct {
	health port.Health
	pubsub port.Pubsub
}

func NewStatusApplication(
	health port.Health,
	pubsub port.Pubsub,
) *StatusApplication {
	return &StatusApplication{
		health: health,
		pubsub: pubsub,
	}
}

func (sa *StatusApplication) HandleRequestDeploymentStatus(d *value.Deployment) {
	releaseName := d.DeploymentName
	health, err := sa.health.CheckPodsHealthy(d.Namespace, releaseName, d.KubeconfigBase64)
	if err != nil {
		log.Printf("failed to check pods health: %v\n", err)
	}

	err = sa.pubsub.PublishDeploymentStatusResponse(&value.HealthStatus{
		DeploymentId: d.DeploymentId,
		Health:       value.Health(health.Health),
		Status:       health.Status,
	})
	if err != nil {
		log.Printf("failed to publish deployment status: %v\n", err)
	}
}
