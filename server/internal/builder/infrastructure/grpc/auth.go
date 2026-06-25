package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"starliner.app/internal/builder/conf"
	"starliner.app/internal/builder/domain/port"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

const defaultResolveRunnerTimeout = 5 * time.Second

type AuthClient struct {
	client  v1.RunnerAuthServiceClient
	timeout time.Duration
}

func NewAuthClient(cfg *conf.Config) (port.AuthClient, error) {
	conn, err := grpc.NewClient(cfg.ApiGrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		client:  v1.NewRunnerAuthServiceClient(conn),
		timeout: defaultResolveRunnerTimeout,
	}, nil
}

func (c *AuthClient) ResolveRunner(ctx context.Context, token string) (*port.ResolvedRunner, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.ResolveRunnerId(ctx, &v1.ResolveRunnerIdRequest{Token: token})
	if err != nil {
		return nil, err
	}

	return &port.ResolvedRunner{
		Id:             resp.GetRunnerId(),
		OrganizationId: resp.OrganizationId,
	}, nil
}

var _ port.AuthClient = (*AuthClient)(nil)
