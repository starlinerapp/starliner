package application

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	teamRepository        interfaces.TeamRepository
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
	teamRepository interfaces.TeamRepository,
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
		teamRepository:        teamRepository,
		normalizerService:     normalizerService,
		organizationService:   organizationService,
		cfg:                   cfg,
	}
}

func (ga *GitHubApplication) VerifySignature(payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(ga.cfg.GithubWebhookSecret))
	mac.Write(payload)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
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

	userTeams, err := ga.teamRepository.GetUserTeams(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	allowedRepoIDs := make(map[int64]bool)

	for _, team := range userTeams {
		teamRepos, err := ga.teamRepository.GetTeamRepositories(ctx, team.Id)
		if err != nil {
			return nil, err
		}
		for _, tr := range teamRepos {
			allowedRepoIDs[tr.GithubRepoId] = true
		}
	}

	var filtered []*value.Repository
	for _, repo := range repos {
		if repo.Id != nil && allowedRepoIDs[*repo.Id] {
			filtered = append(filtered, value.NewRepository(repo))
		}
	}

	return filtered, nil
}

func (ga *GitHubApplication) GetAllRepositories(ctx context.Context, userId int64, organizationId int64) ([]*value.Repository, error) {
	err := ga.organizationService.ValidateUserOrgOwner(ctx, organizationId, userId)
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
	case *value.PushToBranchEvent:
		return ga.triggerBuildsForRepository(ctx, e.RepositoryUrl, e.TargetBranch)
	default:
		return nil
	}
}

func (ga *GitHubApplication) triggerBuildsForRepository(ctx context.Context, repositoryUrl string, branch string) error {
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

		environmentBranch, err := ga.environmentRepository.GetEnvironmentBranch(ctx, deployment.EnvironmentId)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if environmentBranch != branch {
			continue
		}

		b, err := ga.buildRepository.CreateBuild(ctx, deployment.Id, "push")
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
		if ghApp == nil {
			continue
		}

		accessToken, err := ga.gitHub.GetInstallationToken(ctx, ghApp.InstallationID)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		coreArgs := make([]*coreValue.Arg, len(deployment.Args))
		for i, a := range deployment.Args {
			coreArgs[i] = &coreValue.Arg{
				Name:  a.Name,
				Value: a.Value,
			}
		}

		err = ga.queue.PublishBuildTriggered(&coreValue.TriggerBuild{
			BuildId:        b.Id,
			DeploymentId:   deployment.Id,
			ImageName:      fmt.Sprintf("%s/%s", env.Namespace, normalizedServiceName),
			GitUrl:         deployment.GitUrl,
			AccessToken:    accessToken,
			RootDirectory:  deployment.ProjectRepositoryPath,
			DockerfilePath: deployment.DockerfilePath,
			Args:           coreArgs,
		})
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	return errors.Join(errs...)
}

func (ga *GitHubApplication) GetFileContent(ctx context.Context, userId int64, organizationId int64, owner string, repository string, path string) (string, error) {
	err := ga.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return "", err
	}

	ghApp, err := ga.githubAppRepository.GetOrganizationGithubApp(ctx, organizationId)
	if err != nil {
		return "", err
	}

	content, err := ga.gitHub.GetFile(ctx, ghApp.InstallationID, owner, repository, path)
	if err != nil {
		return "", err
	}

	return content, nil
}
