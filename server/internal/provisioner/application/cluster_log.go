package application

import (
	"context"
	"fmt"
	"io"

	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/provisioner/domain/port"
)

type ClusterLogApplication struct {
	streams corePort.KVStore
}

var _ port.LogPublisher = (*ClusterLogApplication)(nil)

func NewClusterLogApplication(
	streams corePort.KVStore,
) *ClusterLogApplication {
	return &ClusterLogApplication{
		streams: streams,
	}
}

func (a *ClusterLogApplication) StreamProvisioningLogs(ctx context.Context, clusterId int64) (io.ReadCloser, error) {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

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
				clusterLogStream(clusterId),
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

func (a *ClusterLogApplication) PublishLogChunk(ctx context.Context, clusterId int64, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	return a.streams.AppendToStream(
		ctx,
		clusterLogStream(clusterId),
		map[string][]byte{
			"data": data,
		},
	)
}

func clusterLogStream(clusterId int64) string {
	return fmt.Sprintf("cluster:%d:logs", clusterId)
}
