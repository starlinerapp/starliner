package application

import (
	"context"
	"fmt"
	"io"

	"starliner.app/internal/cluster/domain/port"
	corePort "starliner.app/internal/core/domain/port"
)

type DeploymentStatusApplication struct {
	deploymentStatus port.DeploymentStatus
	streams          corePort.KVStore
}

var _ port.LogPublisher = (*DeploymentStatusApplication)(nil)

func NewDeploymentStatusApplication(
	deploymentStatus port.DeploymentStatus,
	streams corePort.KVStore,
) *DeploymentStatusApplication {
	return &DeploymentStatusApplication{
		deploymentStatus: deploymentStatus,
		streams:          streams,
	}
}

func (a *DeploymentStatusApplication) StreamDeploymentStatusLogs(
	ctx context.Context,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
	commitHash string,
) (io.ReadCloser, error) {
	return a.deploymentStatus.StreamDeploymentStatusLogs(
		ctx,
		namespace,
		releaseName,
		kubeconfigBase64,
		commitHash,
	)
}

func (a *DeploymentStatusApplication) StreamIngressDeploymentStatusLogs(
	ctx context.Context,
	namespace string,
	releaseName string,
) (io.ReadCloser, error) {
	pr, pw := io.Pipe()

	go func() {
		defer func(pw *io.PipeWriter) {
			_ = pw.Close()
		}(pw)

		lastId := "0"

		for {
			select {
			case <-ctx.Done():
				_ = pw.CloseWithError(ctx.Err())
				return
			default:
			}

			entries, err := a.streams.ReadStream(
				ctx,
				ingressLogStream(namespace, releaseName),
				lastId,
			)
			if err != nil {
				_ = pw.CloseWithError(err)
				return
			}

			for _, entry := range entries {
				lastId = entry.ID

				data, ok := entry.Values["data"]
				if !ok {
					continue
				}

				if _, err := pw.Write(data); err != nil {
					_ = pw.CloseWithError(err)
					return
				}
			}
		}
	}()
	return pr, nil
}

func (a *DeploymentStatusApplication) PublishLogChunk(ctx context.Context, namespace string, releaseName string, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	return a.streams.AppendToStream(
		ctx,
		ingressLogStream(namespace, releaseName),
		map[string][]byte{
			"data": data,
		},
	)
}

func ingressLogStream(namespace string, releaseName string) string {
	return fmt.Sprintf("ingress:%s:%s:logs", namespace, releaseName)
}
