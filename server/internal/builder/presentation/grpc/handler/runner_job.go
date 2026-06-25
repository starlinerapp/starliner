package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"starliner.app/internal/builder/application"
	"starliner.app/internal/core/domain/value"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type RunnerJobHandler struct {
	v1.UnimplementedRunnerJobServiceServer
	runnerJobApplication *application.RunnerJobApplication
}

func NewRunnerJobHandler(runnerJobApplication *application.RunnerJobApplication) *RunnerJobHandler {
	return &RunnerJobHandler{
		runnerJobApplication: runnerJobApplication,
	}
}

func (h *RunnerJobHandler) ClaimJob(
	ctx context.Context,
	_ *v1.ClaimJobRequest,
) (*v1.ClaimJobResponse, error) {
	job, err := h.runnerJobApplication.ClaimJob(ctx, tokenFromContext(ctx))
	if err != nil {
		return nil, err
	}
	if job == nil {
		return &v1.ClaimJobResponse{}, nil
	}

	return &v1.ClaimJobResponse{Job: runnerBuildJobToProto(job)}, nil
}

func (h *RunnerJobHandler) ReportBuildLog(
	ctx context.Context,
	req *v1.ReportBuildLogRequest,
) (*v1.ReportBuildLogResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if err := h.runnerJobApplication.ReportBuildLog(
		ctx,
		tokenFromContext(ctx),
		req.GetBuildId(),
		req.GetData(),
	); err != nil {
		return nil, err
	}

	return &v1.ReportBuildLogResponse{}, nil
}

func (h *RunnerJobHandler) ReportBuildResult(
	ctx context.Context,
	req *v1.ReportBuildResultRequest,
) (*v1.ReportBuildResultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if err := h.runnerJobApplication.ReportBuildResult(
		ctx,
		tokenFromContext(ctx),
		runnerBuildResultFromProto(req.GetResult()),
	); err != nil {
		return nil, err
	}

	return &v1.ReportBuildResultResponse{}, nil
}

func runnerBuildJobToProto(job *value.RunnerBuildJob) *v1.BuildJob {
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

func runnerBuildResultFromProto(result *v1.BuildResult) *value.RunnerBuildResult {
	if result == nil {
		return nil
	}

	buildStatus := value.BuildStatusFailed
	if result.GetStatus() == v1.BuildStatus_BUILD_STATUS_SUCCESS {
		buildStatus = value.BuildStatusSuccess
	}

	return &value.RunnerBuildResult{
		BuildId:      result.GetBuildId(),
		DeploymentId: result.GetDeploymentId(),
		CommitHash:   result.GetCommitHash(),
		Tag:          result.GetTag(),
		ImageName:    result.GetImageName(),
		Logs:         result.GetLogs(),
		Status:       buildStatus,
	}
}
