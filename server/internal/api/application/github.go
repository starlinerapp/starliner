package application

import (
	"context"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
)

type GitHubApplication struct {
	gitHub              port.GitHub
	githubAppRepository interfaces.GithubAppRepository
	organizationService *service.OrganizationService
	cfg                 *conf.Config
}

func NewGitHubApplication(
	gitHub port.GitHub,
	githubAppRepository interfaces.GithubAppRepository,
	organizationService *service.OrganizationService,
	cfg *conf.Config,
) *GitHubApplication {
	return &GitHubApplication{
		gitHub:              gitHub,
		githubAppRepository: githubAppRepository,
		organizationService: organizationService,
		cfg:                 cfg,
	}
}

func (ga *GitHubApplication) GetRepositories(ctx context.Context, userId int64, organizationId int64) ([]*value.Repository, error) {
	err := ga.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	ghApp, err := ga.githubAppRepository.GetOrganizationGithubApp(ctx, organizationId)
	if err != nil {
		return nil, err
	}

	repos, err := ga.gitHub.ListRepositories(ctx, ghApp.InstallationID)
	if err != nil {
		return nil, err
	}

	return value.NewRepositories(repos), nil
}
