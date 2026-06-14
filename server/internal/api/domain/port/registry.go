package port

import "context"

type Registry interface {
	GetRegistryPushToken(ctx context.Context, repository string) (string, error)
}
