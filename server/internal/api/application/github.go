package application

import (
	"context"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/api/domain/value"
)

type GitHubApplication struct {
	gitHub port.GitHub
	cfg    *conf.Config
}

func NewGitHubApplication(
	gitHub port.GitHub,
	cfg *conf.Config,
) *GitHubApplication {
	return &GitHubApplication{
		gitHub: gitHub,
		cfg:    cfg,
	}
}

func (ga *GitHubApplication) GetRepositories(ctx context.Context) ([]*value.Repository, error) {
	repos, err := ga.gitHub.ListRepositories(ctx, 0)
	if err != nil {
		return nil, err
	}

	return value.NewRepositories(repos), nil
}
