package port

import "context"

type AuthClient interface {
	ResolveRunnerId(ctx context.Context, token string) (int64, error)
}
