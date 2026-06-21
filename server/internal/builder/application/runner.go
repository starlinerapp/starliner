package application

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"starliner.app/internal/builder/domain/port"
	corePort "starliner.app/internal/core/domain/port"
)

type RunnerApplication struct {
	authClient  port.AuthClient
	runnerStore corePort.RunnerStore
}

func NewRunnerApplication(authClient port.AuthClient, runnerStore corePort.RunnerStore) *RunnerApplication {
	return &RunnerApplication{
		authClient:  authClient,
		runnerStore: runnerStore,
	}
}

func (a *RunnerApplication) ResolveRunner(ctx context.Context, token string) (*port.ResolvedRunner, error) {
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "missing runner token")
	}

	runner, err := a.authClient.ResolveRunner(ctx, token)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
			return nil, status.Error(codes.Unauthenticated, st.Message())
		}
		return nil, status.Errorf(codes.Unavailable, "resolve runner: %v", err)
	}

	return runner, nil
}

func (a *RunnerApplication) HandleRunnerDeleted(ctx context.Context, runnerId int64) error {
	state, err := a.runnerStore.GetRunner(ctx, runnerId)
	if err != nil {
		return err
	}

	return a.runnerStore.MarkDeleted(ctx, runnerId, state)
}
