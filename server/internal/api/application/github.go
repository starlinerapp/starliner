package application

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"

	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	corePort "starliner.app/internal/core/domain/port"
	coreService "starliner.app/internal/core/domain/service"
)

type GitHubApplication struct {
	gitHub                port.GitHub
	queue                 port.Queue
	registry              port.Registry
	crypto                corePort.Crypto
	deploymentRepository  interfaces.DeploymentRepository
	projectRepository     interfaces.ProjectRepository
	buildRepository       interfaces.BuildRepository
	environmentRepository interfaces.EnvironmentRepository
	githubAppRepository   interfaces.GithubAppRepository
	teamRepository        interfaces.TeamRepository
	parserService         *service.ParserService
	resolverService       *service.ResolverService
	environmentService    *service.EnvironmentService
	gitDeploymentService  *service.GitDeploymentService
	organizationService   *service.OrganizationService
	normalizerService     *coreService.NormalizerService
	cfg                   *conf.Config
}

func NewGitHubApplication(
	gitHub port.GitHub,
	queue port.Queue,
	registry port.Registry,
	crypto corePort.Crypto,
	deploymentRepository interfaces.DeploymentRepository,
	projectRepository interfaces.ProjectRepository,
	buildRepository interfaces.BuildRepository,
	environmentRepository interfaces.EnvironmentRepository,
	githubAppRepository interfaces.GithubAppRepository,
	teamRepository interfaces.TeamRepository,
	parserService *service.ParserService,
	resolverService *service.ResolverService,
	environmentService *service.EnvironmentService,
	gitDeploymentService *service.GitDeploymentService,
	organizationService *service.OrganizationService,
	normalizerService *coreService.NormalizerService,
	cfg *conf.Config,
) *GitHubApplication {
	return &GitHubApplication{
		gitHub:                gitHub,
		queue:                 queue,
		registry:              registry,
		crypto:                crypto,
		deploymentRepository:  deploymentRepository,
		projectRepository:     projectRepository,
		buildRepository:       buildRepository,
		environmentRepository: environmentRepository,
		githubAppRepository:   githubAppRepository,
		teamRepository:        teamRepository,
		parserService:         parserService,
		resolverService:       resolverService,
		environmentService:    environmentService,
		gitDeploymentService:  gitDeploymentService,
		organizationService:   organizationService,
		normalizerService:     normalizerService,
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
	case *value.PullRequestOpenedEvent:
		return ga.createPreviewEnvironment(ctx, e)
	case *value.PullRequestClosedEvent:
		return ga.deletePreviewEnvironment(ctx, e)
	case *value.PushToBranchEvent:
		return ga.triggerBuildsForRepository(ctx, e.RepositoryUrl, e.TargetBranch)
	case *value.GitHubAppInstallationDeletedEvent:
		return ga.deleteGitHubApp(ctx, e)
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
		if deployment.EnvironmentId == nil {
			log.Printf("skipping deployment %d: nil environment id", deployment.Id)
			continue
		}

		environmentBranch, err := ga.environmentRepository.GetEnvironmentBranch(ctx, *deployment.EnvironmentId)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if environmentBranch != branch {
			continue
		}

		if err := ga.gitDeploymentService.RedeployAndBuildForPush(ctx, deployment, branch); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (ga *GitHubApplication) createPreviewEnvironment(ctx context.Context, event *value.PullRequestOpenedEvent) error {
	previewEnv, err := ga.environmentRepository.GetPreviewEnvironment(ctx, event.RepositoryId, event.PrNumber)
	if err != nil {
		return err
	}
	if previewEnv != nil {
		return nil
	}

	productionEnvs, err := ga.projectRepository.GetProjectProductionEnvironmentsByRepositoryUrl(ctx, event.RepositoryUrl)
	if err != nil {
		return err
	}

	var (
		errs        []error
		commentURLs []string
	)

	for _, env := range productionEnvs {
		p, err := ga.environmentRepository.GetEnvironmentProject(ctx, env.Id)
		if err != nil {
			errs = append(errs, err)
		}

		if p.PrEnvironmentsEnabled != nil && !*p.PrEnvironmentsEnabled {
			continue
		}

		if err != nil {
			errs = append(errs, err)
			continue
		}
		randomPrefix := ga.environmentService.RandomPrefix(4)
		previewEnvName := fmt.Sprintf("%s-%s", event.SourceBranch, "preview")
		environmentSlug, err := ga.normalizerService.FormatToDNS1123(previewEnvName)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		namespace, err := ga.normalizerService.FormatToDNS1123(p.Name + "-" + previewEnvName)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		newEnv, err := ga.environmentService.CloneAndProvision(ctx, service.CloneSpec{
			Name:              previewEnvName,
			Namespace:         namespace,
			Slug:              environmentSlug,
			ProjectID:         p.Id,
			SourceEnvironment: env.Id,
			UniquePrefix:      randomPrefix,
			ConnectedBranch:   &event.SourceBranch,
			Preview: &entity.PreviewCloneMetadata{
				GithubRepositoryId: event.RepositoryId,
				PrNumber:           event.PrNumber,
			},
		}, service.ProvisionOptions{
			GitStrategy: service.GitProvisionBuildFromBranch,
			GitBranch:   event.SourceBranch,
			PreviewURLs: &commentURLs,
		})
		if err != nil {
			errs = append(errs, err)
			continue
		}

		ghApp, err := ga.githubAppRepository.GetEnvironmentGithubApp(ctx, newEnv.Id)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if ghApp != nil {
			body, err := buildPreviewEnvironmentComment(commentURLs)
			if err != nil {
				errs = append(errs, err)
			}
			if body != "" {
				err = ga.gitHub.CreatePRComment(
					ctx,
					ghApp.InstallationID,
					event.RepositoryOwner,
					event.RepositoryName,
					event.PrNumber,
					body,
				)
				if err != nil {
					errs = append(errs, err)
				}
			}
		}
	}
	return errors.Join(errs...)
}

func (ga *GitHubApplication) deletePreviewEnvironment(ctx context.Context, event *value.PullRequestClosedEvent) error {
	previewEnv, err := ga.environmentRepository.GetPreviewEnvironment(ctx, event.RepositoryId, event.PrNumber)
	if err != nil {
		return err
	}
	if previewEnv == nil {
		return nil
	}

	env, err := ga.environmentRepository.GetEnvironmentById(ctx, previewEnv.Id)
	if err != nil {
		return err
	}

	if err := ga.environmentService.TearDownEnvironmentDeployments(ctx, env); err != nil {
		return err
	}

	return ga.environmentRepository.DeleteEnvironment(ctx, previewEnv.Id)
}

func (ga *GitHubApplication) deleteGitHubApp(ctx context.Context, event *value.GitHubAppInstallationDeletedEvent) error {
	if event.InstallationId == nil {
		return fmt.Errorf("installation id is nil")
	}
	return ga.githubAppRepository.DeleteGithubApp(ctx, *event.InstallationId)
}

func buildPreviewEnvironmentComment(urls []string) (string, error) {
	seen := make(map[string]struct{})
	unique := make([]string, 0, len(urls))

	for _, u := range urls {
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		unique = append(unique, u)
	}

	if len(unique) == 0 {
		return "", nil
	}

	var sb strings.Builder
	sb.WriteString("Preview links:\n")
	for _, u := range unique {
		_, err := fmt.Fprintf(&sb, "- %s\n", u)
		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
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
