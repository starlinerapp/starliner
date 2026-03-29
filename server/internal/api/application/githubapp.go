package application

import (
	"context"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
)

type GitHubAppApplication struct {
	organizationService *service.OrganizationService
	githubAppRepository interfaces.GithubAppRepository
}

func NewGitHubAppApplication(
	organizationService *service.OrganizationService,
	githubAppRepository interfaces.GithubAppRepository,
) *GitHubAppApplication {
	return &GitHubAppApplication{
		organizationService: organizationService,
		githubAppRepository: githubAppRepository,
	}
}

func (ga *GitHubAppApplication) CreateGitHubApp(ctx context.Context, userId int64, organizationId int64, installationId int64) error {
	err := ga.organizationService.ValidateUserOrgOwner(ctx, organizationId, userId)
	if err != nil {
		return err
	}

	_, err = ga.githubAppRepository.CreateGithubApp(ctx, installationId, organizationId)
	if err != nil {
		return err
	}
	return nil
}
