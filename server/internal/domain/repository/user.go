package repository

import (
	"context"
	"database/sql"
	"errors"
	"starliner.app/internal/domain/entity"
	"starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/infrastructure/postgres/sqlc"
)

type UserRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.UserRepository = (*UserRepository)(nil)

func NewUserRepository(queries *sqlc.Queries) interfaces.UserRepository {
	return &UserRepository{queries: queries}
}

func (u *UserRepository) CreateUser(ctx context.Context, betterAuthID string) (*entity.User, error) {
	user, err := u.queries.CreateUser(ctx, betterAuthID)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		Id:           user.ID,
		BetterAuthId: user.BetterAuthID,
	}, nil
}

func (u *UserRepository) GetUserByBetterAuthId(ctx context.Context, betterAuthID string) (*entity.User, error) {
	user, err := u.queries.GetUserByBetterAuthId(ctx, betterAuthID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity.User{
		Id:           user.ID,
		BetterAuthId: user.BetterAuthID,
	}, nil
}
