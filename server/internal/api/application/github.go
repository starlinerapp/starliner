package application

import (
	"context"
	"errors"
	"fmt"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	coreService "starliner.app/internal/core/domain/service"
	coreValue "starliner.app/internal/core/domain/value"
)

type GitHubApplication struct {
	gitHub                port.GitHub
	queue                 port.Queue
	deploymentRepository  interfaces.DeploymentRepository
	buildRepository       interfaces.BuildRepository
	environmentRepository interfaces.EnvironmentRepository
	githubAppRepository   interfaces.GithubAppRepository
	normalizerService     *coreService.NormalizerService
	organizationService   *service.OrganizationService
	cfg                   *conf.Config
}

func NewGitHubApplication(
	gitHub port.GitHub,
	queue port.Queue,
	deploymentRepository interfaces.DeploymentRepository,
	buildRepository interfaces.BuildRepository,
	environmentRepository interfaces.EnvironmentRepository,
	githubAppRepository interfaces.GithubAppRepository,
	normalizerService *coreService.NormalizerService,
	organizationService *service.OrganizationService,
	cfg *conf.Config,
) *GitHubApplication {
	return &GitHubApplication{
		gitHub:                gitHub,
		queue:                 queue,
		deploymentRepository:  deploymentRepository,
		buildRepository:       buildRepository,
		environmentRepository: environmentRepository,
		githubAppRepository:   githubAppRepository,
		normalizerService:     normalizerService,
		organizationService:   organizationService,
		cfg:                   cfg,
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

func (ga *GitHubApplication) GetRepositoryContents(ctx context.Context, userId int64, organizationId int64, owner string, repository string, repositoryPath string) ([]*value.RepositoryFile, error) {
	err := ga.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	ghApp, err := ga.githubAppRepository.GetOrganizationGithubApp(ctx, organizationId)
	if err != nil {
		return nil, err
	}

	content, err := ga.gitHub.ListRepositoryContents(ctx, ghApp.InstallationID, owner, repository, repositoryPath)
	if err != nil {
		return nil, err
	}

	return value.NewRepositoryFiles(content), nil
}

func (ga *GitHubApplication) HandleGithubWebhook(ctx context.Context, eventType string, payload []byte) error {
	event, err := ga.gitHub.ParseGitEvent(eventType, payload)
	if err != nil {
		return err
	}

	switch e := event.(type) {
	case *value.PullRequestClosedEvent:
		return ga.triggerBuildsForRepository(ctx, e.RepositoryUrl)
	case *value.PushToBranchEvent:
		if e.TargetBranch != "main" {
			return nil
		}
		return ga.triggerBuildsForRepository(ctx, e.RepositoryUrl)
	default:
		return nil
	}
}

func (ga *GitHubApplication) triggerBuildsForRepository(ctx context.Context, repositoryUrl string) error {
	deployments, err := ga.deploymentRepository.GetGitDeploymentsByRepositoryUrl(ctx, repositoryUrl)
	if err != nil {
		return err
	}

	var errs []error

	for _, deployment := range deployments {
		env, err := ga.environmentRepository.GetEnvironmentById(ctx, deployment.EnvironmentId)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		b, err := ga.buildRepository.CreateBuild(ctx, deployment.Id)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		normalizedServiceName, err := ga.normalizerService.FormatToDNS1123(deployment.Name)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		ghApp, err := ga.githubAppRepository.GetEnvironmentGithubApp(ctx, deployment.EnvironmentId)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		accessToken, err := ga.gitHub.GetInstallationToken(ctx, ghApp.InstallationID)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = ga.queue.PublishBuildTriggered(&coreValue.TriggerBuild{
			BuildId:        b.Id,
			DeploymentId:   deployment.Id,
			ImageName:      fmt.Sprintf("%s/%s", env.Namespace, normalizedServiceName),
			GitUrl:         deployment.GitUrl,
			AccessToken:    accessToken,
			RootDirectory:  deployment.ProjectRepositoryPath,
			DockerfilePath: deployment.DockerfilePath,
		})
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	return errors.Join(errs...)
}
