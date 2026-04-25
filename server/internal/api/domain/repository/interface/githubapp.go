package interfaces

import (
	"context"
	"starliner.app/internal/api/domain/entity"
)

type GithubAppRepository interface {
	CreateGithubApp(ctx context.Context, installationID int64, organizationId int64) (*entity.GithubApp, error)
	GetOrganizationGithubApp(ctx context.Context, organizationId int64) (*entity.GithubApp, error)
	GetEnvironmentGithubApp(ctx context.Context, environmentId int64) (*entity.GithubApp, error)
	DeleteGithubApp(ctx context.Context, installationId int64) error
}
