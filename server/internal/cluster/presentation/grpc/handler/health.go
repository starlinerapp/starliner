package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"starliner.app/internal/cluster/application"
	"starliner.app/internal/cluster/infrastructure/k8s"
	"starliner.app/internal/core/domain/value"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type HealthHandler struct {
	v1.UnimplementedHealthServiceServer
	statusApplication *application.StatusApplication
}

func NewHealthHandler(statusApplication *application.StatusApplication) *HealthHandler {
	return &HealthHandler{
		statusApplication: statusApplication,
	}
}

func (h *HealthHandler) GetHealthStatus(
	ctx context.Context,
	req *v1.GetHealthStatusRequest,
) (*v1.GetHealthStatusResponse, error) {
	health, err := h.statusApplication.GetHealthStatus(
		ctx,
		req.GetNamespace(),
		req.GetDeploymentName(),
		req.GetKubeconfigBase64(),
	)
	if err != nil {
		if k8s.IsClusterUnreachable(err) {
			return nil, status.Error(codes.Unavailable, value.ErrClusterUnreachable.Error())
		}
		return nil, status.Errorf(codes.Internal, "check pods health: %v", err)
	}

	return &v1.GetHealthStatusResponse{
		DeploymentId: req.GetDeploymentId(),
		Health:       healthToProto(health.Health),
		Status:       health.Status,
	}, nil
}

func healthToProto(health value.Health) v1.Health {
	switch health {
	case value.Healthy:
		return v1.Health_HEALTH_HEALTHY
	case value.Unhealthy:
		return v1.Health_HEALTH_UNHEALTHY
	default:
		return v1.Health_HEALTH_UNSPECIFIED
	}
}
