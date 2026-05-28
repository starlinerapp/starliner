package port

import coreValue "starliner.app/internal/core/domain/value"

type Pubsub interface {
	PublishDeploymentStatusRequest(deployment *coreValue.Deployment) error
	SubscribeToDeploymentStatusResponse(handler func(health *coreValue.HealthStatus)) error
	SubscribeToReconcileClusterRequest(handler func(request *coreValue.ReconcileClusterRequest)) error
}
