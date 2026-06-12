package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/auth/challenge"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/port"
)

type Client struct {
	cfg        *conf.Config
	httpClient *http.Client
}

type staticCredentialStore struct {
	username string
	password string
}

func (s *staticCredentialStore) Basic(*url.URL) (string, string) {
	return s.username, s.password
}

func (s *staticCredentialStore) RefreshToken(*url.URL, string) string {
	return ""
}

func (s *staticCredentialStore) SetRefreshToken(*url.URL, string, string) {}

func NewClient(cfg *conf.Config) port.Registry {
	return &Client{
		cfg:        cfg,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) GetRepositoryPushToken(ctx context.Context, repository string) (string, error) {
	if repository == "" {
		return "", fmt.Errorf("repository is required")
	}

	baseURL, err := normalizeRegistryURL(c.cfg.ImageRegistryUrl)
	if err != nil {
		return "", err
	}

	challengeManager := challenge.NewSimpleManager()
	pingEndpoints := []string{
		baseURL + "/v2/",
		baseURL + "/v2/" + repository + "/manifests/latest",
	}

	var pingErr error
	for _, endpoint := range pingEndpoints {
		pingErr = c.ping(ctx, challengeManager, endpoint)
		if pingErr == nil {
			break
		}
	}
	if pingErr != nil {
		return "", fmt.Errorf("failed to discover registry auth challenge: %w", pingErr)
	}

	pingURL, err := url.Parse(baseURL + "/v2/")
	if err != nil {
		return "", fmt.Errorf("invalid registry url: %w", err)
	}

	challenges, err := challengeManager.GetChallenges(*pingURL)
	if err != nil {
		return "", err
	}

	creds := &staticCredentialStore{
		username: c.cfg.ImageRegistryUsername,
		password: c.cfg.ImageRegistryPassword,
	}

	transport := c.httpClient.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	tokenHandler := auth.NewTokenHandlerWithOptions(auth.TokenHandlerOptions{
		Transport:   transport,
		Credentials: creds,
		ClientID:    "starliner",
		Scopes: []auth.Scope{
			auth.RepositoryScope{
				Repository: repository,
				Actions:    []string{"pull", "push"},
			},
		},
	})

	for _, ch := range challenges {
		if ch.Scheme != "bearer" {
			continue
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, pingURL.String(), nil)
		if err != nil {
			return "", err
		}

		if err := tokenHandler.AuthorizeRequest(req, ch.Parameters); err != nil {
			return "", fmt.Errorf("failed to authorize registry token request: %w", err)
		}

		authHeader := req.Header.Get("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			return "", fmt.Errorf("registry did not return a bearer token")
		}

		token := strings.TrimPrefix(authHeader, prefix)
		if token == "" {
			return "", fmt.Errorf("registry returned an empty bearer token")
		}

		return token, nil
	}

	return "", fmt.Errorf("registry did not advertise bearer authentication")
}

func (c *Client) ping(ctx context.Context, manager challenge.Manager, endpoint string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("expected 401, got %d", resp.StatusCode)
	}

	return manager.AddResponse(resp)
}

func normalizeRegistryURL(registryURL string) (string, error) {
	registryURL = strings.TrimSpace(registryURL)
	if registryURL == "" {
		return "", fmt.Errorf("registry url is required")
	}

	if !strings.HasPrefix(registryURL, "http://") && !strings.HasPrefix(registryURL, "https://") {
		registryURL = "https://" + registryURL
	}

	parsed, err := url.Parse(registryURL)
	if err != nil {
		return "", err
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("invalid registry url: %s", registryURL)
	}

	parsed.Path = strings.TrimSuffix(parsed.Path, "/")
	parsed.RawQuery = ""
	parsed.Fragment = ""

	return parsed.String(), nil
}
