package repository

import (
	"context"
	"database/sql"
	"errors"
	"starliner.app/internal/domain"
	"starliner.app/internal/infrastructure/db/sqlc"
	interfaces "starliner.app/internal/repository/interface"
)

type UserRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.UserRepository = (*UserRepository)(nil)

func NewUserRepository(queries *sqlc.Queries) interfaces.UserRepository {
	return &UserRepository{queries: queries}
}

func (u *UserRepository) CreateUser(ctx context.Context, betterAuthID string) (*domain.User, error) {
	user, err := u.queries.CreateUser(ctx, betterAuthID)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Id:           user.ID,
		BetterAuthId: user.BetterAuthID,
	}, nil
}

func (u *UserRepository) GetUserByBetterAuthId(ctx context.Context, betterAuthID string) (*domain.User, error) {
	user, err := u.queries.GetUserByBetterAuthId(ctx, betterAuthID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		Id:           user.ID,
		BetterAuthId: user.BetterAuthID,
	}, nil
}
