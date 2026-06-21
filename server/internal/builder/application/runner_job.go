package application

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
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

func (a *RunnerJobApplication) ClaimJob(ctx context.Context, token string) (*v1.BuildJob, error) {
	runner, err := a.runnerApplication.ResolveRunner(ctx, token)
	if err != nil {
		return nil, err
	}

	job, err := a.queue.ClaimRunnerJob(runner.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "claim runner job: %v", err)
	}
	if job == nil {
		return nil, nil
	}

	return toProtoBuildJob(job), nil
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

func (a *RunnerJobApplication) ReportBuildResult(ctx context.Context, token string, result *v1.BuildResult) error {
	if _, err := a.runnerApplication.ResolveRunner(ctx, token); err != nil {
		return err
	}
	if result == nil {
		return status.Error(codes.InvalidArgument, "result is required")
	}
	if result.GetBuildId() <= 0 {
		return status.Error(codes.InvalidArgument, "build id is required")
	}

	if err := a.queue.PublishRunnerJobResult(toValueRunnerBuildResult(result)); err != nil {
		return status.Errorf(codes.Internal, "publish runner job result: %v", err)
	}

	return nil
}

func toProtoBuildJob(job *value.RunnerBuildJob) *v1.BuildJob {
	if job == nil {
		return nil
	}

	protoArgs := make([]*v1.BuildArg, 0, len(job.Args))
	for _, arg := range job.Args {
		if arg == nil {
			continue
		}
		protoArgs = append(protoArgs, &v1.BuildArg{
			Name:  arg.Name,
			Value: arg.Value,
		})
	}

	return &v1.BuildJob{
		BuildId:           job.BuildId,
		DeploymentId:      job.DeploymentId,
		ImageName:         job.ImageName,
		ImageRegistryUrl:  job.ImageRegistryUrl,
		GitUrl:            job.GitUrl,
		BranchName:        job.BranchName,
		AccessToken:       job.AccessToken,
		RegistryPushToken: job.RegistryPushToken,
		RootDirectory:     job.RootDirectory,
		DockerfilePath:    job.DockerfilePath,
		Args:              protoArgs,
	}
}

func toValueRunnerBuildResult(result *v1.BuildResult) *value.RunnerBuildResult {
	status := value.BuildStatusFailed
	if result.GetStatus() == v1.BuildStatus_BUILD_STATUS_SUCCESS {
		status = value.BuildStatusSuccess
	}

	return &value.RunnerBuildResult{
		BuildId:      result.GetBuildId(),
		DeploymentId: result.GetDeploymentId(),
		CommitHash:   result.GetCommitHash(),
		Tag:          result.GetTag(),
		ImageName:    result.GetImageName(),
		Logs:         result.GetLogs(),
		Status:       status,
	}
}
