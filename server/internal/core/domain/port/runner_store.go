package port

import (
	"context"
	"time"

	"starliner.app/internal/core/domain/value"
)

type RunnerRuntimeState struct {
	OrganizationId    *int64
	Status            value.RunnerStatus
	MaxConcurrentJobs int32
	ActiveJobs        int32
}

type RunnerStore interface {
	UpsertHeartbeat(ctx context.Context, runnerId int64, state RunnerRuntimeState, leaseTTL time.Duration) error
	GetRunner(ctx context.Context, runnerId int64) (*RunnerRuntimeState, error)
	IsRunnerDeleted(ctx context.Context, runnerId int64) (bool, error)
	IsRunnerAlive(ctx context.Context, runnerId int64) (bool, error)
	ListOrgRunners(ctx context.Context, orgId int64) ([]int64, error)
	ListGlobalRunners(ctx context.Context) ([]int64, error)
	ListMonitoredRunners(ctx context.Context) ([]int64, error)
	SetRunnerStatus(ctx context.Context, runnerId int64, status value.RunnerStatus) error
	MarkDeleted(ctx context.Context, runnerId int64, state *RunnerRuntimeState) error
}
