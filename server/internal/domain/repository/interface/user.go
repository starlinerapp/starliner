package interfaces

import (
	"context"
	"starliner.app/internal/domain/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, betterAuthID string) (*entity.User, error)
	GetUserByBetterAuthId(ctx context.Context, betterAuthID string) (*entity.User, error)
}
