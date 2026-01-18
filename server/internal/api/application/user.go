package application

import (
	"context"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
)

type UserApplication struct {
	userRepository _interface.UserRepository
}

func NewUserApplication(userRepository _interface.UserRepository) *UserApplication {
	return &UserApplication{userRepository: userRepository}
}

func (us *UserApplication) GetOrCreateUser(ctx context.Context, betterAuthID string) (*value.User, error) {
	user, err := us.userRepository.GetUserByBetterAuthId(ctx, betterAuthID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return value.NewUser(user), nil
	}

	newUser, err := us.userRepository.CreateUser(ctx, betterAuthID)
	if err != nil {
		return nil, err
	}
	return value.NewUser(newUser), nil
}
