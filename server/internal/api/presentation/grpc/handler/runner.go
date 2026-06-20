package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type RunnerHandler struct {
	v1.UnimplementedRunnerAuthServiceServer
	runnerApplication *application.RunnerApplication
}

func NewRunnerHandler(runnerApplication *application.RunnerApplication) *RunnerHandler {
	return &RunnerHandler{
		runnerApplication: runnerApplication,
	}
}

func (h *RunnerHandler) ResolveRunnerId(
	ctx context.Context,
	req *v1.ResolveRunnerIdRequest,
) (*v1.ResolveRunnerIdResponse, error) {
	runnerId, err := h.runnerApplication.ResolveRunnerId(ctx, req.GetToken())
	if err != nil {
		if errors.Is(err, value.ErrInvalidRunnerRegistrationToken) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "resolve runner id: %v", err)
	}

	return &v1.ResolveRunnerIdResponse{
		RunnerId: runnerId,
	}, nil
}
