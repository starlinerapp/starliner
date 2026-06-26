package application

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
)

type RunnerJobApplication struct {
	runnerApplication *RunnerApplication
	queue             port.Queue
	logPublisher      port.LogPublisher
}

func NewRunnerJobApplication(
	runnerApplication *RunnerApplication,
	queue port.Queue,
	logPublisher port.LogPublisher,
) *RunnerJobApplication {
	return &RunnerJobApplication{
		runnerApplication: runnerApplication,
		queue:             queue,
		logPublisher:      logPublisher,
	}
}

func (a *RunnerJobApplication) ClaimJob(ctx context.Context, token string) (*value.RunnerBuildJob, error) {
	runner, err := a.runnerApplication.ResolveRunner(ctx, token)
	if err != nil {
		return nil, err
	}

	job, err := a.queue.ClaimRunnerJob(runner.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "claim runner job: %v", err)
	}

	return job, nil
}

func (a *RunnerJobApplication) ReportBuildLog(ctx context.Context, token string, buildId int64, data []byte) error {
	if _, err := a.runnerApplication.ResolveRunner(ctx, token); err != nil {
		return err
	}
	if buildId <= 0 {
		return status.Error(codes.InvalidArgument, "build id is required")
	}
	if len(data) == 0 {
		return nil
	}

	if a.logPublisher == nil {
		return status.Error(codes.Internal, "log publisher not configured")
	}

	if err := a.logPublisher.PublishLogChunk(buildId, data); err != nil {
		return status.Errorf(codes.Internal, "publish build log chunk: %v", err)
	}

	return nil
}

func (a *RunnerJobApplication) ReportBuildResult(ctx context.Context, token string, result *value.RunnerBuildResult) error {
	if _, err := a.runnerApplication.ResolveRunner(ctx, token); err != nil {
		return err
	}
	if result == nil {
		return status.Error(codes.InvalidArgument, "result is required")
	}
	if result.BuildId <= 0 {
		return status.Error(codes.InvalidArgument, "build id is required")
	}

	if err := a.queue.PublishRunnerJobResult(result); err != nil {
		return status.Errorf(codes.Internal, "publish runner job result: %v", err)
	}

	return nil
}
