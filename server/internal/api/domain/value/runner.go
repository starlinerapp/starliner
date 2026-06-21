package value

import (
	"errors"
	"time"

	"starliner.app/internal/api/domain/entity"
)

var ErrInvalidRunnerRegistrationToken = errors.New("invalid runner registration token")
var ErrRunnerNotFound = errors.New("runner not found")

type CreateRunnerResult struct {
	Id        int64
	Token     string
	ExpiresAt time.Time
}

type Runner struct {
	Id                int64
	OrganizationId    *int64
	IsGlobal          bool
	Name              *string
	Status            string
	Labels            []string
	MaxConcurrentJobs int32
	DisabledAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewRunner(r *entity.Runner) *Runner {
	return &Runner{
		Id:                r.Id,
		OrganizationId:    r.OrganizationId,
		IsGlobal:          r.OrganizationId == nil,
		Name:              r.Name,
		Status:            r.Status,
		Labels:            r.Labels,
		MaxConcurrentJobs: r.MaxConcurrentJobs,
		DisabledAt:        r.DisabledAt,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

func NewRunners(rs []*entity.Runner) []*Runner {
	runners := make([]*Runner, len(rs))
	for i, r := range rs {
		runners[i] = NewRunner(r)
	}
	return runners
}
