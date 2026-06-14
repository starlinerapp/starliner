package response

import (
	"time"

	"starliner.app/internal/api/domain/value"
)

type CreateRunner struct {
	Id        int64     `json:"id" binding:"required"`
	Token     string    `json:"token" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
}

func NewCreateRunner(result *value.CreateRunnerResult) CreateRunner {
	return CreateRunner{
		Id:        result.Id,
		Token:     result.Token,
		ExpiresAt: result.ExpiresAt,
	}
}
