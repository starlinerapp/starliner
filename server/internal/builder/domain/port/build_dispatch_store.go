package port

import "context"

type DispatchedBuild struct {
	BuildId      int64
	DeploymentId int64
}

type BuildDispatchStore interface {
	Register(ctx context.Context, runnerId int64, buildId int64, deploymentId int64) error
	Unregister(ctx context.Context, buildId int64) error
	ListByRunner(ctx context.Context, runnerId int64) ([]DispatchedBuild, error)
}
