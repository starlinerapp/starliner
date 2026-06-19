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
}
