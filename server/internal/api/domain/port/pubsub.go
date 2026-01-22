package port

import coreValue "starliner.app/internal/core/domain/value"

type Pubsub interface {
	PublishDeploymentStatusRequest(deployment coreValue.Deployment) error
}
