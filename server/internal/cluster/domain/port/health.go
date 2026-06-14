package port

import (
	"context"

	"starliner.app/internal/cluster/domain/value"
)

type Health interface {
	CheckPodsHealthy(ctx context.Context, namespace string, releaseName string, kubeconfigBase64 string) (*value.HealthStatus, error)
}
