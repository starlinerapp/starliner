package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"starliner.app/internal/builder/application"
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

	return &v1.ClaimJobResponse{Job: job}, nil
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
		req.GetResult(),
	); err != nil {
		return nil, err
	}

	return &v1.ReportBuildResultResponse{}, nil
}
