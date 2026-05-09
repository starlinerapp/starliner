package port

import "context"

type AuthClient interface {
	Login(ctx context.Context, email string, password string) error
	Logout(ctx context.Context) error
	LoadToken() (string, error)
}
