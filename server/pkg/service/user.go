package service

import (
	"context"
	"starliner.app/pkg/domain"
	"starliner.app/pkg/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
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
