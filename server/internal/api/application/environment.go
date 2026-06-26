package application

import (
	"context"
	"fmt"
	"sort"

	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	coreService "starliner.app/internal/core/domain/service"
)

type EnvironmentApplication struct {
	organizationService   *service.OrganizationService
	environmentService    *service.EnvironmentService
	normalizerService     *coreService.NormalizerService
	buildRepository       interfaces.BuildRepository
	environmentRepository interfaces.EnvironmentRepository
	projectRepository     interfaces.ProjectRepository
}

func NewEnvironmentApplication(
	organizationService *service.OrganizationService,
	environmentService *service.EnvironmentService,
	normalizerService *coreService.NormalizerService,
	buildRepository interfaces.BuildRepository,
	environmentRepository interfaces.EnvironmentRepository,
	projectRepository interfaces.ProjectRepository,
) *EnvironmentApplication {
	return &EnvironmentApplication{
		organizationService:   organizationService,
		environmentService:    environmentService,
		normalizerService:     normalizerService,
		buildRepository:       buildRepository,
		environmentRepository: environmentRepository,
		projectRepository:     projectRepository,
	}
}

func (ea *EnvironmentApplication) CreateEnvironment(
	ctx context.Context,
	name string,
	userId int64,
	organizationId int64,
	projectId int64,
	sourceEnvironmentId *int64,
) (*value.Environment, error) {
	err := ea.organizationService.ValidateUserInOrg(ctx, organizationId, userId)
	if err != nil {
		return nil, err
	}

	environmentSlug, err := ea.normalizerService.FormatToDNS1123(name)
	if err != nil {
		return nil, err
	}

	project, err := ea.projectRepository.GetProject(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}

	namespace, err := ea.normalizerService.FormatToDNS1123(project.Name + "-" + name)
	if err != nil {
		return nil, err
	}

	if sourceEnvironmentId != nil {
		randomPrefix := ea.environmentService.RandomPrefix(4)
		env, err := ea.environmentService.CloneAndProvision(ctx, service.CloneSpec{
			Name:              name,
			Namespace:         namespace,
			Slug:              environmentSlug,
			ProjectID:         projectId,
			SourceEnvironment: *sourceEnvironmentId,
			UniquePrefix:      randomPrefix,
		}, service.ProvisionOptions{
			GitStrategy:            service.GitProvisionReuseSourceImage,
			SourceEnvIDForGitReuse: *sourceEnvironmentId,
		})
		if err != nil {
			return nil, err
		}
		return value.NewEnvironment(env), nil
	}

	env, err := ea.environmentRepository.CreateEnvironment(ctx, name, namespace, environmentSlug, projectId)
	if err != nil {
		return nil, err
	}

	return value.NewEnvironment(env), nil
}

func (ea *EnvironmentApplication) GetEnvironmentDeployments(ctx context.Context, environmentId int64, userId int64) (*value.Deployments, error) {
	if err := ea.environmentService.ValidateUserPermission(ctx, userId, environmentId); err != nil {
		return nil, err
	}

	return ea.environmentService.ListEnvironmentDeployments(ctx, environmentId)
}

func (ea *EnvironmentApplication) GetEnvironmentGitDeploymentBuilds(ctx context.Context, userId int64, environmentId int64) ([]*value.GitDeploymentBuild, error) {
	err := ea.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return nil, err
	}

	builds, err := ea.environmentRepository.GetEnvironmentGitDeploymentBuilds(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	ingressBuilds, err := ea.environmentRepository.GetEnvironmentIngressDeploymentBuilds(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	imageBuilds, err := ea.environmentRepository.GetEnvironmentImageDeploymentBuilds(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	databaseBuilds, err := ea.environmentRepository.GetEnvironmentDatabaseDeploymentBuilds(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	allBuilds := append(builds, ingressBuilds...)
	allBuilds = append(allBuilds, imageBuilds...)
	allBuilds = append(allBuilds, databaseBuilds...)
	sort.Slice(allBuilds, func(i, j int) bool {
		return allBuilds[i].CreatedAt.After(allBuilds[j].CreatedAt)
	})

	valueBuilds := make([]*value.GitDeploymentBuild, len(allBuilds))
	for i, b := range allBuilds {
		valueBuilds[i] = &value.GitDeploymentBuild{
			BuildId:                 b.BuildId,
			DeploymentId:            b.DeploymentId,
			DeploymentName:          b.DeploymentName,
			DeploymentRolloutStatus: b.DeploymentRolloutStatus,
			CommitHash:              b.CommitHash,
			Source:                  b.Source,
			Status:                  value.BuildStatus(b.Status),
			CreatedAt:               b.CreatedAt,
		}
	}
	return valueBuilds, nil
}

func (ea *EnvironmentApplication) GetEnvironmentBranch(ctx context.Context, userId int64, environmentId int64) (string, error) {
	err := ea.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return "", err
	}

	branch, err := ea.environmentRepository.GetEnvironmentBranch(ctx, environmentId)
	if err != nil {
		return "", err
	}
	return branch, nil
}

func (ea *EnvironmentApplication) UpdateEnvironmentBranch(ctx context.Context, userId int64, environmentId int64, branch string) error {
	err := ea.environmentService.ValidateUserPermission(ctx, userId, environmentId)
	if err != nil {
		return err
	}

	err = ea.environmentRepository.UpdateEnvironmentBranch(ctx, environmentId, branch)
	if err != nil {
		return err
	}
	return nil
}

func (ea *EnvironmentApplication) DeleteEnvironment(ctx context.Context, userId int64, environmentId int64) error {
	if err := ea.environmentService.ValidateUserPermission(ctx, userId, environmentId); err != nil {
		return err
	}
	env, err := ea.environmentRepository.GetEnvironmentById(ctx, environmentId)
	if err != nil {
		return err
	}

	if env.Slug == value.EnvironmentProductionSlug {
		return fmt.Errorf("cannot delete production environment")
	}

	if err := ea.environmentService.TearDownEnvironmentDeployments(ctx, env); err != nil {
		return err
	}
	return ea.environmentRepository.DeleteEnvironment(ctx, environmentId)
}
