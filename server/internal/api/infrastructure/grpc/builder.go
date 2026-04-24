package grpc

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/port"
	v2 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type BuilderClient struct {
	buildLogClient v2.BuildLogServiceClient
}

func NewBuilderClient(cfg *conf.Config) (port.BuilderClient, error) {
	conn, err := grpc.NewClient(cfg.BuilderGrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &BuilderClient{
		buildLogClient: v2.NewBuildLogServiceClient(conn),
	}, nil
}

func (c *BuilderClient) StreamBuildLogs(ctx context.Context, buildId int64, w io.Writer) error {
	stream, err := c.buildLogClient.StreamBuildLogs(ctx, &v2.StreamBuildLogRequest{BuildId: buildId})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		_, err = w.Write(resp.Chunk)
		if err != nil {
			return err
		}
	}
}
