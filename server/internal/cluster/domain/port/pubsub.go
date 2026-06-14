package port

import "starliner.app/internal/core/domain/value"

type Pubsub interface {
	PublishReconcileClusterRequest(request *value.ReconcileClusterRequest) error
}
