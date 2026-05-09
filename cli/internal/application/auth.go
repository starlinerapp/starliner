package application

import (
	"context"

	"starliner.app/cli/internal/domain/port"
)

type AuthApplication struct {
	authClient port.AuthClient
}

func NewAuthApplication(
	authClient port.AuthClient,
) *AuthApplication {
	return &AuthApplication{
		authClient: authClient,
	}
}

func (a *AuthApplication) Login(ctx context.Context, email string, password string) error {
	return a.authClient.Login(ctx, email, password)
}

func (a *AuthApplication) Logout(ctx context.Context) error {
	return a.authClient.Logout(ctx)
}
