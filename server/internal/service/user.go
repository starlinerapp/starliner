package service

import (
	"context"
	interfaces "starliner.app/internal/repository/interface"
	"starliner.app/internal/service/model"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(userRepository interfaces.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (us *UserService) GetOrCreateUser(ctx context.Context, betterAuthID string) (*model.User, error) {
	user, err := us.userRepository.GetUserByBetterAuthId(ctx, betterAuthID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return &model.User{
			Id:           user.Id,
			BetterAuthId: user.BetterAuthId,
		}, nil
	}

	newUser, err := us.userRepository.CreateUser(ctx, betterAuthID)
	if err != nil {
		return nil, err
	}
	return &model.User{
		Id:           newUser.Id,
		BetterAuthId: newUser.BetterAuthId,
	}, nil
}
