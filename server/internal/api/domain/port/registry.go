package port

import "context"

type Registry interface {
	GetRepositoryPushToken(ctx context.Context, repository string) (string, error)
}
