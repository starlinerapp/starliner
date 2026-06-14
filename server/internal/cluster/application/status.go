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
	ctx context.Context,
	namespace string,
	deploymentName string,
	kubeconfigBase64 string,
) (*value.HealthStatus, error) {
	health, err := sa.health.CheckPodsHealthy(ctx, namespace, deploymentName, kubeconfigBase64)
	if err != nil {
		return nil, err
	}

	return &value.HealthStatus{
		Health: value.Health(health.Health),
		Status: health.Status,
	}, nil
}
