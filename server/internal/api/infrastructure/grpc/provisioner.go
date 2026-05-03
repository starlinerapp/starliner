package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/port"
	v2 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type ProvisionerClient struct {
	clusterTtyServiceClient v2.ClusterTTYServiceClient
	provisioningLogsClient  v2.ProvisioningLogServiceClient
}

func NewProvisionerClient(cfg *conf.Config) (port.ProvisionerClient, error) {
	conn, err := grpc.NewClient(cfg.ProvisionerGrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &ProvisionerClient{
		clusterTtyServiceClient: v2.NewClusterTTYServiceClient(conn),
		provisioningLogsClient:  v2.NewProvisioningLogServiceClient(conn),
	}, nil
}

func (c *ProvisionerClient) StreamProvisioningLogs(ctx context.Context, clusterId int64, w io.Writer) error {
	stream, err := c.provisioningLogsClient.StreamProvisioningLogs(ctx, &v2.StreamProvisioningLogRequest{
		ClusterId: clusterId,
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

		_, err = w.Write(resp.Chunk)
		if err != nil {
			return err
		}
	}
}

func (c *ProvisionerClient) OpenTTY(
	ctx context.Context,
	ip string,
	user string,
	pemKey []byte,
	stdin io.Reader,
	stdout io.Writer,
	sizes <-chan port.TerminalSize,
) error {
	stream, err := c.clusterTtyServiceClient.OpenTTY(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&v2.OpenClusterTTYRequest{
		Payload: &v2.OpenClusterTTYRequest_Session{
			Session: &v2.ClusterTTYSession{
				Ip:     ip,
				User:   user,
				PemKey: pemKey,
			},
		},
	})
	if err != nil {
		return err
	}

	errCh := make(chan error, 3)

	// Forward stdout from server to writer
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}
			if _, err := stdout.Write(msg.Stdout); err != nil {
				errCh <- err
				return
			}
		}
	}()

	// Forward stdin from reader to server
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stdin.Read(buf)
			if n > 0 {
				if err := stream.Send(&v2.OpenClusterTTYRequest{
					Payload: &v2.OpenClusterTTYRequest_Stdin{
						Stdin: buf[:n],
					},
				}); err != nil {
					errCh <- err
					return
				}
			}
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	// Forward terminal resize events
	go func() {
		for {
			select {
			case size, ok := <-sizes:
				if !ok {
					return
				}
				if err := stream.Send(&v2.OpenClusterTTYRequest{
					Payload: &v2.OpenClusterTTYRequest_Size{
						Size: &v2.ClusterTerminalSize{
							Cols: uint32(size.Columns),
							Rows: uint32(size.Rows),
						},
					},
				}); err != nil {
					errCh <- err
					return
				}
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}
	}()

	err = <-errCh
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}
