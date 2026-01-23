package port

import "starliner.app/internal/cluster/domain/value"

type Health interface {
	CheckPodsHealthy(releaseName string, kubeconfigPath string) (*value.HealthStatus, error)
}
