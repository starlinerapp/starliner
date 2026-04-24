package port

import (
	"context"

	coreValue "starliner.app/internal/core/domain/value"
)

type Pubsub interface {
	PublishDeploymentStatusRequest(deployment *coreValue.Deployment) error
	SubscribeToDeploymentStatusResponse(handler func(health *coreValue.HealthStatus)) error

	// SubscribeToBuildLogs opens a transient NATS subscription for the given
	// build and delivers log chunks on the returned channel. The channel is
	// closed when ctx is canceled or when the builder publishes a chunk with
	// End=true. The returned cancel releases the subscription early.
	SubscribeToBuildLogs(ctx context.Context, buildId int64) (<-chan *coreValue.BuildLogChunk, func(), error)
}
