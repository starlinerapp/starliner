package interfaces

import (
	"context"
	"starliner.app/pkg/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, betterAuthID string) (*domain.User, error)
	GetUserByBetterAuthId(ctx context.Context, betterAuthID string) (*domain.User, error)
}
