package service

import (
	"context"
	"fmt"

	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
	coreService "starliner.app/internal/core/domain/service"
	coreValue "starliner.app/internal/core/domain/value"
)

type RedeployGitConfig struct {
	Port                  string
	ProjectRepositoryPath string
	DockerfilePath        string
	Envs                  []*value.EnvVar
	Args                  []*value.Arg
}

type GitDeploymentService struct {
	deploymentRepository  interfaces.DeploymentRepository
	buildRepository       interfaces.BuildRepository
	environmentRepository interfaces.EnvironmentRepository
	githubAppRepository   interfaces.GithubAppRepository
	gitHub                port.GitHub
	queue                 port.Queue
	registry              port.Registry
	normalizer            *coreService.NormalizerService
}

func NewGitDeploymentService(
	deploymentRepository interfaces.DeploymentRepository,
	buildRepository interfaces.BuildRepository,
	environmentRepository interfaces.EnvironmentRepository,
	githubAppRepository interfaces.GithubAppRepository,
	gitHub port.GitHub,
	queue port.Queue,
	registry port.Registry,
	normalizer *coreService.NormalizerService,
) *GitDeploymentService {
	return &GitDeploymentService{
		deploymentRepository:  deploymentRepository,
		buildRepository:       buildRepository,
		environmentRepository: environmentRepository,
		githubAppRepository:   githubAppRepository,
		gitHub:                gitHub,
		queue:                 queue,
		registry:              registry,
		normalizer:            normalizer,
	}
}

func (s *GitDeploymentService) Redeploy(
	ctx context.Context,
	existing *entity.GitDeployment,
	config RedeployGitConfig,
) (*entity.GitDeployment, error) {
	if existing.EnvironmentId == nil {
		return nil, fmt.Errorf("deployment %d has nil environment id", existing.Id)
	}

	if err := s.deploymentRepository.SoftDeleteDeployment(ctx, existing.Id); err != nil {
		return nil, err
	}

	newDeployment, err := s.deploymentRepository.CreateGitDeployment(
		ctx,
		*existing.EnvironmentId,
		existing.Name,
		config.Port,
		existing.GitUrl,
		config.ProjectRepositoryPath,
		config.DockerfilePath,
		config.Envs,
		config.Args,
	)
	if err != nil {
		return nil, err
	}

	if err := s.deploymentRepository.RepointIngressPathsTargetDeployment(ctx, existing.Id, newDeployment.Id); err != nil {
		return nil, err
	}

	return newDeployment, nil
}

func (s *GitDeploymentService) TriggerBuild(
	ctx context.Context,
	deployment *entity.GitDeployment,
	branch string,
	buildSource string,
	args []*value.Arg,
) error {
	if deployment.EnvironmentId == nil {
		return fmt.Errorf("deployment %d has nil environment id", deployment.Id)
	}
	environmentID := *deployment.EnvironmentId

	env, err := s.environmentRepository.GetEnvironmentById(ctx, environmentID)
	if err != nil {
		return err
	}

	normalizedServiceName, err := s.normalizer.FormatToDNS1123(deployment.Name)
	if err != nil {
		return err
	}

	organization, err := s.environmentRepository.GetEnvironmentOrganization(ctx, environmentID)
	if err != nil {
		return err
	}

	imageName := fmt.Sprintf("%s/%s/%s", organization.Slug, env.Namespace, normalizedServiceName)

	registryPushToken, err := s.registry.GetRegistryPushToken(ctx, imageName)
	if err != nil {
		return err
	}

	b, err := s.buildRepository.CreateBuild(ctx, deployment.Id, buildSource)
	if err != nil {
		return err
	}

	ghApp, err := s.githubAppRepository.GetEnvironmentGithubApp(ctx, environmentID)
	if err != nil {
		return err
	}
	if ghApp == nil {
		return nil
	}

	accessToken, err := s.gitHub.GetInstallationToken(ctx, ghApp.InstallationID)
	if err != nil {
		return err
	}

	coreArgs := make([]*coreValue.Arg, len(args))
	for i, a := range args {
		coreArgs[i] = &coreValue.Arg{
			Name:  a.Name,
			Value: a.Value,
		}
	}

	return s.queue.PublishBuildTriggered(&coreValue.TriggerBuild{
		BuildId:           b.Id,
		DeploymentId:      deployment.Id,
		OrganizationId:    organization.Id,
		ImageName:         imageName,
		GitUrl:            deployment.GitUrl,
		BranchName:        branch,
		AccessToken:       accessToken,
		RegistryPushToken: registryPushToken,
		RootDirectory:     deployment.ProjectRepositoryPath,
		DockerfilePath:    deployment.DockerfilePath,
		Args:              coreArgs,
	})
}

func (s *GitDeploymentService) RedeployAndBuildForPush(
	ctx context.Context,
	existing *entity.GitDeployment,
	branch string,
) error {
	envs := entityEnvVarsToValue(existing.EnvVars)
	args := entityArgsToValue(existing.Args)

	newDeployment, err := s.Redeploy(ctx, existing, RedeployGitConfig{
		Port:                  existing.Port,
		ProjectRepositoryPath: existing.ProjectRepositoryPath,
		DockerfilePath:        existing.DockerfilePath,
		Envs:                  envs,
		Args:                  args,
	})
	if err != nil {
		return err
	}

	return s.TriggerBuild(ctx, newDeployment, branch, "push", args)
}

func entityEnvVarsToValue(envs []*entity.EnvVar) []*value.EnvVar {
	result := make([]*value.EnvVar, len(envs))
	for i, e := range envs {
		result[i] = &value.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		}
	}
	return result
}

func entityArgsToValue(args []*entity.Arg) []*value.Arg {
	result := make([]*value.Arg, len(args))
	for i, a := range args {
		result[i] = &value.Arg{
			Name:  a.Name,
			Value: a.Value,
		}
	}
	return result
}
