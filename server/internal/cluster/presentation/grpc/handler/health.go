package handler

import (
	"context"

	"starliner.app/internal/cluster/application"
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
	deployment := &value.Deployment{
		DeploymentId:     req.GetDeploymentId(),
		Namespace:        req.GetNamespace(),
		DeploymentName:   req.GetDeploymentName(),
		KubeconfigBase64: req.GetKubeconfigBase64(),
		ClusterId:        req.GetClusterId(),
		OrganizationId:   req.GetOrganizationId(),
		ProvisioningId:   req.GetProvisioningId(),
	}

	health, err := h.statusApplication.GetHealthStatus(ctx, deployment)
	if err != nil {
		return nil, err
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
