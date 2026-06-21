package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"starliner.app/internal/builder/conf"
	"starliner.app/internal/builder/domain/port"
	v1 "starliner.app/internal/core/infrastructure/grpc/proto/v1"
)

type AuthClient struct {
	client v1.RunnerAuthServiceClient
}

func NewAuthClient(cfg *conf.Config) (port.AuthClient, error) {
	conn, err := grpc.NewClient(cfg.ApiGrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		client: v1.NewRunnerAuthServiceClient(conn),
	}, nil
}

func (c *AuthClient) ResolveRunner(ctx context.Context, token string) (*port.ResolvedRunner, error) {
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
