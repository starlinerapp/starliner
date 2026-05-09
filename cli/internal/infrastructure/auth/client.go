package auth

import (
	"starliner.app/cli/internal/conf"
	openapi "starliner.app/cli/internal/infrastructure/auth/generated/client"
)

type Client struct {
	apiAuthClient *openapi.APIClient
	conf          *conf.Config
}

func NewClient(
	conf *conf.Config,
) *Client {
	apiAuthClient := openapi.NewAPIClient(&openapi.Configuration{
		Host: conf.AuthBaseUrl,
	})
	return &Client{
		apiAuthClient: apiAuthClient,
	}
}

func (c *Client) Login() error {
	return nil
}
