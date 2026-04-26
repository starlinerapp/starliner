package application

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	corePort "starliner.app/internal/core/domain/port"
	coreService "starliner.app/internal/core/domain/service"
	coreValue "starliner.app/internal/core/domain/value"
	"strconv"
	"strings"
)

type GitHubApplication struct {
	gitHub                port.GitHub
	queue                 port.Queue
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
	organizationService   *service.OrganizationService
	normalizerService     *coreService.NormalizerService
	cfg                   *conf.Config
}

func NewGitHubApplication(
	gitHub port.GitHub,
	queue port.Queue,
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
	organizationService *service.OrganizationService,
	normalizerService *coreService.NormalizerService,
	cfg *conf.Config,
) *GitHubApplication {
	return &GitHubApplication{
		gitHub:                gitHub,
		queue:                 queue,
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
			errs = append(errs, fmt.Errorf("deployment %d has nil environment id", deployment.Id))
			continue
		}
		environmentID := *deployment.EnvironmentId

		env, err := ga.environmentRepository.GetEnvironmentById(ctx, environmentID)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		environmentBranch, err := ga.environmentRepository.GetEnvironmentBranch(ctx, environmentID)
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

		ghApp, err := ga.githubAppRepository.GetEnvironmentGithubApp(ctx, environmentID)
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
			BranchName:     branch,
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
		newEnv, err := ga.environmentRepository.CreatePreviewEnvironment(
			ctx,
			previewEnvName,
			namespace,
			environmentSlug,
			p.Id,
			env.Id,
			randomPrefix,
			&event.SourceBranch,
			event.RepositoryId,
			event.PrNumber,
		)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		deployments, err := ga.getEnvironmentDeployments(ctx, newEnv.Id)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		cluster, err := ga.environmentRepository.GetEnvironmentCluster(ctx, env.Id)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		kubeconfigBase64, err := ga.crypto.Decrypt(*cluster.Kubeconfig)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		ingressDeployments := deployments.Ingresses
		for _, d := range ingressDeployments {
			coreHosts := make([]coreValue.IngressHost, 0, len(d.IngressHosts))
			for _, h := range d.IngressHosts {
				if h.Host != "" {
					commentURLs = append(commentURLs, "http://"+h.Host)
				}

				ch := coreValue.IngressHost{
					Host: h.Host,
				}
				ch.Paths = make([]coreValue.IngressPath, 0, len(h.Paths))

				for _, p := range h.Paths {
					target, err := ga.environmentRepository.GetEnvironmentDeploymentByName(ctx, p.ServiceName, newEnv.Id)
					if err != nil {
						errs = append(errs, err)
						continue
					}

					targetPort, err := strconv.Atoi(target.Port)
					if err != nil {
						errs = append(errs, err)
						continue
					}

					normalizedServiceName, err := ga.normalizerService.FormatToDNS1123(p.ServiceName)
					if err != nil {
						errs = append(errs, err)
						continue
					}

					ch.Paths = append(ch.Paths, coreValue.IngressPath{
						Path:        p.Path,
						PathType:    coreValue.PathType(p.PathType),
						ServiceName: normalizedServiceName,
						ServicePort: targetPort,
					})
				}
				coreHosts = append(coreHosts, ch)
			}
			err = ga.queue.PublishDeployIngress(&coreValue.IngressDeployment{
				IngressHosts:     coreHosts,
				DeploymentId:     d.Id,
				DeploymentName:   d.ServiceName,
				Namespace:        newEnv.Namespace,
				KubeconfigBase64: kubeconfigBase64,
			})
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}

		databaseDeployments := deployments.Databases
		for _, d := range databaseDeployments {
			normalizedServiceName, err := ga.normalizerService.FormatToDNS1123(d.ServiceName)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			err = ga.queue.PublishDeployDatabase(&coreValue.Deployment{
				Namespace:        newEnv.Namespace,
				DeploymentId:     d.Id,
				DeploymentName:   normalizedServiceName,
				KubeconfigBase64: kubeconfigBase64,
			})
			if err != nil {
				log.Printf("error publishing: %v", err)
			}
		}

		imageDeployments := deployments.Images
		for _, d := range imageDeployments {
			deploymentPort, err := strconv.Atoi(d.Port)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			coreEnvs := value.ToCoreEnvVars(d.EnvVars)

			normalizedDeploymentName, err := ga.normalizerService.FormatToDNS1123(d.ServiceName)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			err = ga.queue.PublishDeployImage(&coreValue.ImageDeployment{
				DeploymentId:     d.Id,
				DeploymentName:   normalizedDeploymentName,
				Namespace:        newEnv.Namespace,
				KubeconfigBase64: kubeconfigBase64,
				ImageName:        d.ImageName,
				ImageTag:         d.Tag,
				Port:             deploymentPort,
				VolumeSizeMiB:    d.VolumeSizeMiB,
				VolumeMountPath:  d.VolumeMountPath,
				EnvVars:          coreEnvs,
			})
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}

		gitDeployments := deployments.GitDeployments
		for _, d := range gitDeployments {
			b, err := ga.buildRepository.CreateBuild(ctx, d.Id, "push")
			if err != nil {
				errs = append(errs, err)
				continue
			}

			normalizedServiceName, err := ga.normalizerService.FormatToDNS1123(d.ServiceName)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			ghApp, err := ga.githubAppRepository.GetEnvironmentGithubApp(ctx, newEnv.Id)
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

			coreArgs := make([]*coreValue.Arg, len(d.Args))
			for i, a := range d.Args {
				coreArgs[i] = &coreValue.Arg{
					Name:  a.Name,
					Value: a.Value,
				}
			}

			err = ga.queue.PublishBuildTriggered(&coreValue.TriggerBuild{
				BuildId:        b.Id,
				DeploymentId:   d.Id,
				ImageName:      fmt.Sprintf("%s/%s", newEnv.Namespace, normalizedServiceName),
				GitUrl:         d.GitUrl,
				BranchName:     event.SourceBranch,
				AccessToken:    accessToken,
				RootDirectory:  d.ProjectRepositoryPath,
				DockerfilePath: d.DockerfilePath,
				Args:           coreArgs,
			})
			if err != nil {
				errs = append(errs, err)
				continue
			}
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

	deployments, err := ga.getEnvironmentDeployments(ctx, previewEnv.Id)
	if err != nil {
		return err
	}

	type deploymentEntry struct {
		id          int64
		serviceName string
	}

	var allDeployments []deploymentEntry
	for _, d := range deployments.Ingresses {
		allDeployments = append(allDeployments, deploymentEntry{d.Id, d.ServiceName})
	}
	for _, d := range deployments.Images {
		allDeployments = append(allDeployments, deploymentEntry{d.Id, d.ServiceName})
	}
	for _, d := range deployments.GitDeployments {
		allDeployments = append(allDeployments, deploymentEntry{d.Id, d.ServiceName})
	}
	for _, d := range deployments.Databases {
		allDeployments = append(allDeployments, deploymentEntry{d.Id, d.ServiceName})
	}

	for _, d := range allDeployments {
		cluster, err := ga.deploymentRepository.GetDeploymentCluster(ctx, d.id)
		if err != nil {
			return err
		}

		env, err := ga.environmentRepository.GetEnvironmentById(ctx, previewEnv.Id)
		if err != nil {
			return err
		}

		if cluster.Kubeconfig == nil {
			return fmt.Errorf("cluster kubeconfig is nil")
		}
		kubeconfigBase64, err := ga.crypto.Decrypt(*cluster.Kubeconfig)
		if err != nil {
			return err
		}

		normalizedDeploymentName, err := ga.normalizerService.FormatToDNS1123(d.serviceName)
		if err != nil {
			return err
		}

		if err = ga.deploymentRepository.SoftDeleteDeploymentVolume(ctx, d.id); err != nil {
			return err
		}

		if err = ga.queue.PublishDeleteDeployment(&coreValue.Deployment{
			DeploymentId:     d.id,
			DeploymentName:   normalizedDeploymentName,
			Namespace:        env.Namespace,
			KubeconfigBase64: kubeconfigBase64,
		}); err != nil {
			log.Printf("error publishing: %v", err)
		}
	}
	return ga.environmentRepository.DeleteEnvironment(ctx, previewEnv.Id)
}

func (ga *GitHubApplication) deleteGitHubApp(ctx context.Context, event *value.GitHubAppInstallationDeletedEvent) error {
	if event.InstallationId == nil {
		return fmt.Errorf("installation id is nil")
	}
	return ga.githubAppRepository.DeleteGithubApp(ctx, *event.InstallationId)
}

func (ga *GitHubApplication) getEnvironmentDeployments(ctx context.Context, environmentId int64) (*value.Deployments, error) {
	ingresses, err := ga.environmentRepository.GetEnvironmentIngressDeployments(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	git, err := ga.environmentRepository.GetEnvironmentGitDeployments(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	gitDeployments := make([]*value.GitDeployment, len(git))
	for i, d := range git {
		normalizedServiceName, err := ga.normalizerService.FormatToDNS1123(d.Name)
		if err != nil {
			return nil, err
		}

		internalEndpoint := fmt.Sprintf("%s:%s", normalizedServiceName, d.Port)
		gitDeployments[i] = value.NewGitDeployment(d, internalEndpoint)
	}

	images, err := ga.environmentRepository.GetEnvironmentImageDeployments(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	imageDeployments := make([]*value.ImageDeployment, len(images))
	for i, d := range images {
		normalizedServiceName, err := ga.normalizerService.FormatToDNS1123(d.ServiceName)
		if err != nil {
			return nil, err
		}

		internalEndpoint := fmt.Sprintf("%s:%s", normalizedServiceName, d.Port)
		imageDeployments[i] = value.NewImageDeployment(d, internalEndpoint)
	}

	databases, err := ga.environmentRepository.GetEnvironmentDatabaseDeployments(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	databaseDeployments := make([]*value.DatabaseDeployment, len(databases))
	for i, d := range databases {
		var password *string

		if d.Password != nil {
			decrypted, err := ga.crypto.Decrypt(*d.Password)
			if err != nil {
				return nil, err
			}
			password = &decrypted
		}

		normalizedServiceName, err := ga.normalizerService.FormatToDNS1123(d.ServiceName)
		if err != nil {
			return nil, err
		}

		internalEndpoint := fmt.Sprintf("%s-rw:%s", normalizedServiceName, d.Port)
		databaseDeployments[i] = &value.DatabaseDeployment{
			Id:               d.Id,
			ServiceName:      d.ServiceName,
			InternalEndpoint: internalEndpoint,
			Status:           d.Status,
			Database:         d.Database,
			Username:         d.Username,
			Password:         password,
			Port:             d.Port,
		}
	}

	return &value.Deployments{
		Databases:      databaseDeployments,
		Images:         imageDeployments,
		Ingresses:      value.NewIngressDeployments(ingresses),
		GitDeployments: gitDeployments,
	}, nil
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
