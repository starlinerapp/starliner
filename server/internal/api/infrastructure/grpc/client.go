package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"starliner.app/internal/api/domain/port"
	v2 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type Client struct {
	client v2.LogsServiceClient
}

func NewClient() (port.GrpcClient, error) {
	conn, err := grpc.NewClient("server-cluster:57400", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: v2.NewLogsServiceClient(conn),
	}, nil
}

func (c *Client) StreamLogs(
	ctx context.Context,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
) error {
	stream, err := c.client.StreamLogs(ctx, &v2.StreamLogsRequest{
		Namespace:        namespace,
		ReleaseName:      releaseName,
		KubeconfigBase64: kubeconfigBase64,
	})
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
		fmt.Print(string(resp.Chunk))
	}
}
