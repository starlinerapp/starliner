package auth

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"starliner.app/cli/internal/domain/port"
	openapi "starliner.app/cli/internal/infrastructure/auth/generated/client"
)

type Client struct {
	apiAuthClient *openapi.APIClient
}

func NewClient() port.AuthClient {
	apiAuthClient := openapi.NewAPIClient(
		openapi.NewConfiguration())
	return &Client{
		apiAuthClient: apiAuthClient,
	}
}

func (c *Client) Login(ctx context.Context, email string, password string) error {
	resp, httpResp, err := c.apiAuthClient.DefaultAPI.
		SignInEmail(ctx).
		SignInEmailRequest(openapi.SignInEmailRequest{
			Email:    email,
			Password: password,
		}).
		Execute()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", loginError(err, httpResp))
	}

	return saveToken(resp.Token)
}

func (c *Client) Logout(ctx context.Context) error {
	_, _, err := c.apiAuthClient.DefaultAPI.
		SignOut(ctx).
		Execute()
	if err != nil {
		return err
	}

	return deleteToken()
}

func (c *Client) LoadToken() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(filepath.Join(dir, "starliner", "credentials"))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func saveToken(token string) error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(dir, "starliner")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(configDir, "credentials"), []byte(token), 0600)
}

func deleteToken() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join(dir, "starliner", "credentials"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
