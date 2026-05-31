package handler

import (
	"io"
	"log"

	"google.golang.org/grpc"
	"starliner.app/internal/cluster/application"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type DeploymentStatusLogHandler struct {
	v1.UnimplementedDeploymentStatusLogServiceServer
	v1.UnimplementedIngressDeploymentStatusLogServiceServer
	deploymentStatusApplication *application.DeploymentStatusApplication
}

func NewDeploymentStatusLogHandler(
	deploymentStatusApplication *application.DeploymentStatusApplication,
) *DeploymentStatusLogHandler {
	return &DeploymentStatusLogHandler{
		deploymentStatusApplication: deploymentStatusApplication,
	}
}

func (h *DeploymentStatusLogHandler) StreamDeploymentStatusLogs(
	req *v1.StreamDeploymentStatusLogsRequest,
	stream grpc.ServerStreamingServer[v1.StreamDeploymentStatusLogsResponse],
) error {
	rc, err := h.deploymentStatusApplication.StreamDeploymentStatusLogs(
		stream.Context(),
		req.GetNamespace(),
		req.GetReleaseName(),
		req.GetKubeconfigBase64(),
		req.GetCommitHash(),
	)
	if err != nil {
		return err
	}
	defer func(rc io.ReadCloser) {
		if err := rc.Close(); err != nil {
			log.Printf("failed to close deployment status log reader: %v", err)
		}
	}(rc)

	buf := make([]byte, 64*1024)
	for {
		n, readErr := rc.Read(buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])
			if err := stream.Send(&v1.StreamDeploymentStatusLogsResponse{Chunk: chunk}); err != nil {
				return err
			}
		}
		if readErr == io.EOF {
			return nil
		}
		if readErr != nil {
			return readErr
		}
	}
}

func (h *DeploymentStatusLogHandler) StreamIngressDeploymentStatusLogs(
	req *v1.StreamIngressDeploymentStatusLogsRequest,
	stream grpc.ServerStreamingServer[v1.StreamIngressDeploymentStatusLogsResponse],
) error {
	rc, err := h.deploymentStatusApplication.StreamIngressDeploymentStatusLogs(
		stream.Context(),
		req.GetNamespace(),
		req.GetReleaseName(),
	)
	if err != nil {
		return err
	}
	defer func(rc io.ReadCloser) {
		if err := rc.Close(); err != nil {
			log.Printf("failed to close ingress deployment status log reader: %v", err)
		}
	}(rc)

	buf := make([]byte, 64*1024)
	for {
		n, readErr := rc.Read(buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])
			if err := stream.Send(&v1.StreamIngressDeploymentStatusLogsResponse{Chunk: chunk}); err != nil {
				return err
			}
			if readErr == io.EOF {
				return nil
			}
			if readErr != nil {
				return readErr
			}
		}
	}
}
