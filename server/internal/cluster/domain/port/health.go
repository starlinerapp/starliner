package port

import "starliner.app/internal/cluster/domain/value"

type Health interface {
	CheckPodsHealthy(namespace string, releaseName string, kubeconfigBase64 string) (*value.HealthStatus, error)
}
