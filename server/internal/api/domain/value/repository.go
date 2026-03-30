package value

import (
	"starliner.app/internal/api/domain/port"
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

func NewRepository(r *port.Repository) *Repository {
	return &Repository{
		Id:          r.Id,
		Name:        r.Name,
		FullName:    r.FullName,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		PushedAt:    r.PushedAt,
		UpdatedAt:   r.UpdatedAt,
		CloneURL:    r.CloneURL,
	}
}

func NewRepositories(rs []*port.Repository) []*Repository {
	repositories := make([]*Repository, len(rs))
	for i, r := range rs {
		repositories[i] = NewRepository(r)
	}
	return repositories
}
