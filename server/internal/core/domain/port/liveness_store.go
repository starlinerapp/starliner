package port

import (
	"context"
	"time"
)

type LivenessStore interface {
	MarkAlive(ctx context.Context, key string, ttl time.Duration) error
	IsAlive(ctx context.Context, key string) (bool, error)
	AddMonitoredRunner(ctx context.Context, runnerId int64) error
	ListMonitoredRunners(ctx context.Context) ([]int64, error)
	RemoveMonitoredRunner(ctx context.Context, runnerId int64) error
	IsRunnerDeleted(ctx context.Context, runnerId int64) (bool, error)
	GetRunnerStatus(ctx context.Context, runnerId int64) (string, error)
	SetRunnerStatus(ctx context.Context, runnerId int64, status string) error
}
