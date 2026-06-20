package interfaces

import (
	"context"
	"time"

	"starliner.app/internal/api/domain/entity"
)

type RunnerRepository interface {
	CreateRunnerWithRegistrationToken(
		ctx context.Context,
		organizationId int64,
		tokenHash string,
		expiresAt time.Time,
	) (*entity.Runner, error)
	RegisterRunner(
		ctx context.Context,
		tokenHash string,
		name string,
		labels []string,
		maxConcurrentJobs int32,
	) error
	GetOrganizationRunners(ctx context.Context, organizationId int64) ([]*entity.Runner, error)
	GetRunnerIdByRegistrationToken(ctx context.Context, tokenHash string) (int64, error)
	UpdateRunnerStatus(ctx context.Context, runnerId int64, status string) error
	GetRunnerByOrganization(ctx context.Context, runnerId int64, organizationId int64) (*entity.Runner, error)
	DeleteRunner(ctx context.Context, runnerId int64, organizationId int64) error
}
