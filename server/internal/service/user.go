package service

import (
	"context"
	"starliner.app/internal/domain"
	interfaces "starliner.app/internal/repository/interface"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(userRepository interfaces.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) GetOrCreateUser(ctx context.Context, betterAuthID string) (*domain.User, error) {
	user, err := us.userRepository.GetUserByBetterAuthId(ctx, betterAuthID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return us.userRepository.CreateUser(ctx, betterAuthID)
	}
	return user, nil
}
