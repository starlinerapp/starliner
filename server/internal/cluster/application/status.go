package application

import (
	"context"

	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
)

type StatusApplication struct {
	health port.Health
}

func NewStatusApplication(health port.Health) *StatusApplication {
	return &StatusApplication{
		health: health,
	}
}

func (sa *StatusApplication) GetHealthStatus(
	_ context.Context,
	d *value.Deployment,
) (*value.HealthStatus, error) {
	health, err := sa.health.CheckPodsHealthy(d.Namespace, d.DeploymentName, d.KubeconfigBase64)
	if err != nil {
		return nil, err
	}

	return &value.HealthStatus{
		Health: value.Health(health.Health),
		Status: health.Status,
	}, nil
}
