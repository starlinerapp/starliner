package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"starliner.app/cli/internal/conf"
	openapi "starliner.app/cli/internal/infrastructure/auth/generated/client"
)

type Client struct {
	apiAuthClient *openapi.APIClient
	conf          *conf.Config
}

func NewClient() *Client {
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

func loginError(err error, httpResp *http.Response) error {
	openAPIErr, ok := errors.AsType[*openapi.GenericOpenAPIError](err)
	if !ok || httpResp == nil {
		return err
	}

	switch m := openAPIErr.Model().(type) {
	case openapi.SocialSignIn400Response:
		return fmt.Errorf("status %d: %s", httpResp.StatusCode, m.Message)
	case openapi.SocialSignIn403Response:
		if m.Message != nil {
			return fmt.Errorf("status %d: %s", httpResp.StatusCode, *m.Message)
		}
	}

	if body := openAPIErr.Body(); len(body) > 0 {
		return fmt.Errorf("status %d: %s", httpResp.StatusCode, body)
	}
	return fmt.Errorf("status %d", httpResp.StatusCode)
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

	return os.Remove(filepath.Join(dir, "starliner", "credentials"))
}
