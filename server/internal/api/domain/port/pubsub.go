package port

import (
	"context"

	coreValue "starliner.app/internal/core/domain/value"
)

type Pubsub interface {
	PublishDeploymentStatusRequest(deployment *coreValue.Deployment) error
	SubscribeToDeploymentStatusResponse(handler func(health *coreValue.HealthStatus)) error
	SubscribeToBuildLogs(ctx context.Context, buildId int64) (<-chan *coreValue.BuildLogChunk, func(), error)
}
