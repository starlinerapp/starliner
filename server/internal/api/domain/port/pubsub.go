package port

import coreValue "starliner.app/internal/core/domain/value"

type Pubsub interface {
	SubscribeToReconcileClusterRequest(handler func(request *coreValue.ReconcileClusterRequest)) error
}
