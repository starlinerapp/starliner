package port

import (
	"context"
	"time"
)

type Repository struct {
	Id          *int64
	Name        *string
	FullName    *string
	Description *string
	CreatedAt   *time.Time
	PushedAt    *time.Time
	UpdatedAt   *time.Time
	CloneURL    *string
}

type RepositoryFile struct {
	Name *string
	Path *string
	Type *string
	SHA  *string
	Size *int
	URL  string
}

type GitHub interface {
	ListRepositories(ctx context.Context, installationId int64) ([]*Repository, error)
	GetInstallationToken(ctx context.Context, installationId int64) (string, error)
	ListRepositoryContents(
		ctx context.Context,
		installationId int64,
		owner string,
		repository string,
		path string,
	) ([]*RepositoryFile, error)
}
