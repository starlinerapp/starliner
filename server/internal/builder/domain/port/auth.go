package port

import "context"

type ResolvedRunner struct {
	Id             int64
	OrganizationId int64
}

type AuthClient interface {
	ResolveRunner(ctx context.Context, token string) (*ResolvedRunner, error)
}
