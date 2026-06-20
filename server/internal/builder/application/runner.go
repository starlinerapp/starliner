package application

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"starliner.app/internal/builder/domain/port"
)

type RunnerApplication struct {
	runnerAuthClient port.RunnerAuthClient
}

func NewRunnerApplication(runnerAuthClient port.RunnerAuthClient) *RunnerApplication {
	return &RunnerApplication{
		runnerAuthClient: runnerAuthClient,
	}
}

func (a *RunnerApplication) ResolveRunnerId(ctx context.Context, token string) (int64, error) {
	if token == "" {
		return 0, status.Error(codes.Unauthenticated, "missing runner token")
	}

	runnerId, err := a.runnerAuthClient.ResolveRunnerId(ctx, token)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
			return 0, status.Error(codes.Unauthenticated, st.Message())
		}
		return 0, status.Errorf(codes.Unavailable, "resolve runner id: %v", err)
	}

	return runnerId, nil
}
