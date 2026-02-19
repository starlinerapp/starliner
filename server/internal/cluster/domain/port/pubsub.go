package port

import "starliner.app/internal/core/domain/value"

type Pubsub interface {
	SubscribeToDeploymentStatusRequest(func(deployment *value.DatabaseDeployment)) error
	PublishDeploymentStatusResponse(health *value.HealthStatus) error
}
