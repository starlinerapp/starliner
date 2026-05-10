package api

import (
	"starliner.app/cli/internal/domain/port"
	openapi "starliner.app/cli/internal/infrastructure/api/generated/client"
)

type Client struct {
	apiClient *openapi.APIClient
}

func NewClient() port.APIClient {
	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{
		{
			URL: "https://dev.starliner.app",
		},
	}

	apiClient := openapi.NewAPIClient(cfg)
	return &Client{
		apiClient: apiClient,
	}
}
