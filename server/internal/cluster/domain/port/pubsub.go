package port

import "starliner.app/internal/core/domain/value"

type Pubsub interface {
	SubscribeToDeploymentStatusRequest(func(deployment *value.Deployment)) error
	PublishDeploymentStatusResponse(health *value.HealthStatus) error
}
