package response

import (
	"starliner.app/internal/api/domain/value"
	"time"
)

type Repository struct {
	Id          *int64     `json:"id" binding:"required"`
	Name        *string    `json:"name" binding:"required"`
	FullName    *string    `json:"full_name" binding:"required"`
	Owner       *string    `json:"owner" binding:"required"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at" binding:"required"`
	PushedAt    *time.Time `json:"pushed_at" binding:"required"`
	UpdatedAt   *time.Time `json:"updated_at" binding:"required"`
	CloneURL    *string    `json:"clone_url" binding:"required"`
}

func NewRepository(repo *value.Repository) Repository {
	return Repository{
		Id:          repo.Id,
		Name:        repo.Name,
		FullName:    repo.FullName,
		Owner:       repo.Owner,
		Description: repo.Description,
		CreatedAt:   repo.CreatedAt,
		PushedAt:    repo.PushedAt,
		UpdatedAt:   repo.UpdatedAt,
		CloneURL:    repo.CloneURL,
	}
}

func NewRepositories(repos []*value.Repository) []Repository {
	res := make([]Repository, len(repos))
	for i, r := range repos {
		res[i] = NewRepository(r)
	}
	return res
}
