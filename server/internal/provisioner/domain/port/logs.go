package port

import "context"

type LogPublisher interface {
	PublishLogChunk(ctx context.Context, clusterId int64, data []byte) error
}
