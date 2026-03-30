package application

import (
	"context"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
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

func (ga *GitHubAppApplication) GetGithubApp(ctx context.Context, userId int64, organizationId int64) (*value.GithubApp, error) {
	err := ga.organizationService.ValidateUserOrgOwner(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	ghApp, err := ga.githubAppRepository.GetOrganizationGithubApp(ctx, organizationId)
	if err != nil {
		return nil, err
	}

	if ghApp == nil {
		return nil, nil
	}

	return &value.GithubApp{
		InstallationID: ghApp.InstallationID,
		OrganizationID: ghApp.OrganizationID,
		CreatedAt:      ghApp.CreatedAt,
	}, nil
}
