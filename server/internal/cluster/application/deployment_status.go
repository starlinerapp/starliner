package application

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"starliner.app/internal/cluster/domain/port"
	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

type DeploymentStatusApplication struct {
	deploymentStatus port.DeploymentStatus
	streams          corePort.KVStore
	queue            port.Queue
	activePolls      sync.Map
}

var _ port.LogPublisher = (*DeploymentStatusApplication)(nil)

func NewDeploymentStatusApplication(
	deploymentStatus port.DeploymentStatus,
	streams corePort.KVStore,
	queue port.Queue,
) *DeploymentStatusApplication {
	return &DeploymentStatusApplication{
		deploymentStatus: deploymentStatus,
		streams:          streams,
		queue:            queue,
	}
}

func (a *DeploymentStatusApplication) StreamDeploymentStatusLogs(
	ctx context.Context,
	deploymentId int64,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
	commitHash string,
) (io.ReadCloser, error) {
	a.ensureWorkloadPoll(deploymentId, namespace, releaseName, kubeconfigBase64, commitHash)
	return a.streamLogs(ctx, workloadLogStream(deploymentId, namespace, releaseName))
}

func (a *DeploymentStatusApplication) StreamIngressDeploymentStatusLogs(
	ctx context.Context,
	deploymentId int64,
	namespace string,
	releaseName string,
) (io.ReadCloser, error) {
	return a.streamLogs(ctx, ingressLogStream(deploymentId, namespace, releaseName))
}

func (a *DeploymentStatusApplication) ensureWorkloadPoll(
	deploymentId int64,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
	commitHash string,
) {
	if _, loaded := a.activePolls.LoadOrStore(deploymentId, struct{}{}); loaded {
		return
	}

	go func() {
		defer a.activePolls.Delete(deploymentId)

		ctx := context.Background()
		rc, err := a.deploymentStatus.StreamDeploymentStatusLogs(
			ctx,
			namespace,
			releaseName,
			kubeconfigBase64,
			commitHash,
		)
		if err != nil {
			log.Printf("failed to stream workload deployment status logs: %v", err)
			return
		}
		defer func() {
			_ = rc.Close()
		}()

		var logBuf strings.Builder
		buf := make([]byte, 32*1024)
		for {
			n, readErr := rc.Read(buf)
			if n > 0 {
				chunk := buf[:n]
				logBuf.Write(chunk)
				if err := a.publishLogChunk(ctx, workloadLogStream(deploymentId, namespace, releaseName), chunk); err != nil {
					log.Printf("failed to publish workload log chunk: %v", err)
				}
			}
			if readErr == io.EOF {
				if err := a.queue.PublishDeploymentStatusLogsCompleted(&value.DeploymentStatusLogsCompleted{
					DeploymentId: deploymentId,
					Logs:         logBuf.String(),
				}); err != nil {
					log.Printf("failed to publish deployment status logs completed: %v", err)
				}
				return
			}
			if readErr != nil {
				log.Printf("failed to read workload deployment status logs: %v", readErr)
				return
			}
		}
	}()
}

func (a *DeploymentStatusApplication) streamLogs(ctx context.Context, streamName string) (io.ReadCloser, error) {
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

			entries, err := a.streams.ReadStream(ctx, streamName, lastId)
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

func (a *DeploymentStatusApplication) PublishLogChunk(ctx context.Context, deploymentId int64, namespace string, releaseName string, data []byte) error {
	return a.publishLogChunk(ctx, ingressLogStream(deploymentId, namespace, releaseName), data)
}

func (a *DeploymentStatusApplication) publishLogChunk(ctx context.Context, streamName string, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	return a.streams.AppendToStream(
		ctx,
		streamName,
		map[string][]byte{
			"data": data,
		},
	)
}

func ingressLogStream(deploymentId int64, namespace string, releaseName string) string {
	return fmt.Sprintf("ingress:%d:%s:%s:logs", deploymentId, namespace, releaseName)
}

func workloadLogStream(deploymentId int64, namespace string, releaseName string) string {
	return fmt.Sprintf("workload:%d:%s:%s:logs", deploymentId, namespace, releaseName)
}
