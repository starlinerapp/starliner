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

type Runner struct {
	Id                int64      `json:"id" binding:"required"`
	OrganizationId    *int64     `json:"organization_id"`
	IsGlobal          bool       `json:"is_global"`
	Name              *string    `json:"name"`
	Status            string     `json:"status" binding:"required"`
	Labels            []string   `json:"labels" binding:"required"`
	MaxConcurrentJobs int32      `json:"max_concurrent_jobs" binding:"required"`
	DisabledAt        *time.Time `json:"disabled_at"`
	CreatedAt         time.Time  `json:"created_at" binding:"required"`
	UpdatedAt         time.Time  `json:"updated_at" binding:"required"`
}

func NewRunner(runner *value.Runner) Runner {
	return Runner{
		Id:                runner.Id,
		OrganizationId:    runner.OrganizationId,
		IsGlobal:          runner.IsGlobal,
		Name:              runner.Name,
		Status:            runner.Status,
		Labels:            runner.Labels,
		MaxConcurrentJobs: runner.MaxConcurrentJobs,
		DisabledAt:        runner.DisabledAt,
		CreatedAt:         runner.CreatedAt,
		UpdatedAt:         runner.UpdatedAt,
	}
}

func NewRunners(runners []*value.Runner) []Runner {
	res := make([]Runner, len(runners))
	for i, runner := range runners {
		res[i] = NewRunner(runner)
	}
	return res
}
