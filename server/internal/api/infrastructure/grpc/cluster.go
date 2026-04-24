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

type ClusterClient struct {
	logsServiceClient v2.LogsServiceClient
	ttyServiceClient  v2.TTYServiceClient
}

func NewClusterClient(cfg *conf.Config) (port.ClusterClient, error) {
	conn, err := grpc.NewClient(cfg.ClusterGrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &ClusterClient{
		logsServiceClient: v2.NewLogsServiceClient(conn),
		ttyServiceClient:  v2.NewTTYServiceClient(conn),
	}, nil
}

func (c *ClusterClient) StreamLogs(
	ctx context.Context,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
	w io.Writer,
) error {
	stream, err := c.logsServiceClient.StreamLogs(ctx, &v2.StreamLogsRequest{
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

		_, err = w.Write(resp.Chunk)
		if err != nil {
			return err
		}
	}
}

func (c *ClusterClient) OpenTTY(
	ctx context.Context,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
	stdin io.Reader,
	stdout io.Writer,
	sizes <-chan port.TerminalSize,
) error {
	stream, err := c.ttyServiceClient.OpenTTY(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&v2.OpenTTYRequest{
		Payload: &v2.OpenTTYRequest_Session{
			Session: &v2.TTYSession{
				Namespace:        namespace,
				ReleaseName:      releaseName,
				KubeconfigBase64: kubeconfigBase64,
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
				if err := stream.Send(&v2.OpenTTYRequest{
					Payload: &v2.OpenTTYRequest_Stdin{
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
				if err := stream.Send(&v2.OpenTTYRequest{
					Payload: &v2.OpenTTYRequest_Size{
						Size: &v2.TerminalSize{
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
