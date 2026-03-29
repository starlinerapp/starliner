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

type GitHub interface {
	ListRepositories(ctx context.Context, installationId int64) ([]*Repository, error)
}
