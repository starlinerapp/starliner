package application

import (
	"context"
	"starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/domain/value"
)

type UserApplication struct {
	userRepository interfaces.UserRepository
}

func NewUserApplication(userRepository interfaces.UserRepository) *UserApplication {
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
