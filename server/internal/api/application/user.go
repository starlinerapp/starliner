package application

import (
	"context"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
)

type UserApplication struct {
	email          port.Email
	userRepository interfaces.UserRepository
}

func NewUserApplication(email port.Email, userRepository interfaces.UserRepository) *UserApplication {
	return &UserApplication{email: email, userRepository: userRepository}
}

func (us *UserApplication) SendVerificationEmail(ctx context.Context, to string, verificationURL string) error {
	return us.email.SendVerificationEmail(to, port.VerifyData{
		VerificationLink: verificationURL,
	})
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
