package application

import (
	"context"
	"fmt"
	"log"
	"strings"

	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/domain/value"
	corePort "starliner.app/internal/core/domain/port"
	coreService "starliner.app/internal/core/domain/service"
	coreValue "starliner.app/internal/core/domain/value"
)

type ProjectApplication struct {
	normalizerService     *coreService.NormalizerService
	organizationService   *service.OrganizationService
	teamService           *service.TeamService
	projectRepository     interfaces.ProjectRepository
	environmentRepository interfaces.EnvironmentRepository
	deploymentRepository  interfaces.DeploymentRepository
	queue                 port.Queue
	crypto                corePort.Crypto
}

func NewProjectApplication(
	normalizerService *coreService.NormalizerService,
	organizationService *service.OrganizationService,
	teamService *service.TeamService,
	projectRepository interfaces.ProjectRepository,
	environmentRepository interfaces.EnvironmentRepository,
	deploymentRepository interfaces.DeploymentRepository,
	queue port.Queue,
	crypto corePort.Crypto,
) *ProjectApplication {
	return &ProjectApplication{
		normalizerService:     normalizerService,
		organizationService:   organizationService,
		teamService:           teamService,
		projectRepository:     projectRepository,
		environmentRepository: environmentRepository,
		deploymentRepository:  deploymentRepository,
		queue:                 queue,
		crypto:                crypto,
	}
}

func (pa *ProjectApplication) CreateProject(ctx context.Context, name string, organizationId int64, clusterId int64, userId int64, teamId int64) (*value.Project, error) {
	err := pa.teamService.ValidateUserAndClusterInTeam(ctx, userId, teamId, clusterId)
	if err != nil {
		return nil, err
	}

	productionEnvName := "Production"
	namespace, err := pa.normalizerService.FormatToDNS1123(name + "-" + productionEnvName)
	if err != nil {
		return nil, err
	}

	project, err := pa.projectRepository.CreateProjectWithEnvironment(
		ctx,
		name,
		namespace,
		productionEnvName,
		strings.ToLower(productionEnvName),
		teamId,
		clusterId,
	)
	if err != nil {
		return nil, err
	}

	return value.NewProject(project), nil
}

func (pa *ProjectApplication) GetProject(ctx context.Context, projectId int64, userId int64) (*value.Project, error) {
	project, err := pa.projectRepository.GetProject(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}
	return value.NewProject(project), nil
}

func (pa *ProjectApplication) DeleteProject(ctx context.Context, projectId int64, userId int64) error {
	envs, err := pa.projectRepository.GetProjectEnvironments(ctx, projectId, userId)
	if err != nil {
		return err
	}
	for _, env := range envs {
		if err := pa.deleteEnvironmentDeploymentsFromCluster(ctx, env); err != nil {
			return err
		}
		if err := pa.environmentRepository.DeleteEnvironment(ctx, env.Id); err != nil {
			return err
		}
	}
	return pa.projectRepository.DeleteProject(ctx, projectId, userId)
}

func (pa *ProjectApplication) deleteEnvironmentDeploymentsFromCluster(
	ctx context.Context,
	env *entity.Environment,
) error {
	ingresses, err := pa.environmentRepository.GetEnvironmentIngressDeployments(ctx, env.Id)
	if err != nil {
		return err
	}

	gitDeployments, err := pa.environmentRepository.GetEnvironmentGitDeployments(ctx, env.Id)
	if err != nil {
		return err
	}

	images, err := pa.environmentRepository.GetEnvironmentImageDeployments(ctx, env.Id)
	if err != nil {
		return err
	}

	databases, err := pa.environmentRepository.GetEnvironmentDatabaseDeployments(ctx, env.Id)
	if err != nil {
		return err
	}

	type projectDeploymentDelete struct {
		id          int64
		serviceName string
	}
	var deletions []projectDeploymentDelete

	for _, d := range ingresses {
		deletions = append(deletions, projectDeploymentDelete{d.Id, d.Name})
	}

	for _, d := range gitDeployments {
		deletions = append(deletions, projectDeploymentDelete{d.Id, d.Name})
	}

	for _, d := range images {
		deletions = append(deletions, projectDeploymentDelete{d.Id, d.ServiceName})
	}

	for _, d := range databases {
		deletions = append(deletions, projectDeploymentDelete{d.Id, d.ServiceName})
	}

	for _, d := range deletions {
		cluster, err := pa.deploymentRepository.GetDeploymentCluster(ctx, d.id)
		if err != nil {
			return err
		}

		if cluster.Kubeconfig == nil {
			return fmt.Errorf("cluster kubeconfig is nil")
		}
		kubeconfigBase64, err := pa.crypto.Decrypt(*cluster.Kubeconfig)
		if err != nil {
			return err
		}

		normalizedDeploymentName, err := pa.normalizerService.FormatToDNS1123(d.serviceName)
		if err != nil {
			return err
		}

		if err = pa.deploymentRepository.SoftDeleteDeploymentVolume(ctx, d.id); err != nil {
			return err
		}

		if err = pa.queue.PublishDeleteDeployment(&coreValue.Deployment{
			DeploymentId:     d.id,
			DeploymentName:   normalizedDeploymentName,
			Namespace:        env.Namespace,
			KubeconfigBase64: kubeconfigBase64,
		}); err != nil {
			log.Printf("error publishing: %v", err)
		}
	}
	return nil
}

func (pa *ProjectApplication) GetProjectCluster(ctx context.Context, projectId int64, userId int64) (*value.ProjectCluster, error) {
	cluster, err := pa.projectRepository.GetProjectCluster(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}
	return value.NewProjectCluster(cluster), nil
}

func (pa *ProjectApplication) GetProjectEnvironments(ctx context.Context, projectId int64, userId int64) ([]*value.Environment, error) {
	environments, err := pa.projectRepository.GetProjectEnvironments(ctx, projectId, userId)
	if err != nil {
		return nil, err
	}

	valueEnvironments := make([]*value.Environment, len(environments))
	for i, e := range environments {
		valueEnvironments[i] = &value.Environment{
			Id:   e.Id,
			Name: e.Name,
			Slug: e.Slug,
		}
	}
	return valueEnvironments, nil
}

func (pa *ProjectApplication) GetProjectPreviewEnvironmentEnabled(ctx context.Context, projectId int64, userId int64) (bool, error) {
	enabled, err := pa.projectRepository.GetProjectPreviewEnvironmentEnabled(ctx, projectId, userId)
	if err != nil {
		return false, err
	}
	return enabled, nil
}

func (pa *ProjectApplication) ToggleProjectPreviewEnvironmentEnabled(ctx context.Context, projectId int64, userId int64) (bool, error) {
	enabled, err := pa.projectRepository.ToggleProjectPreviewEnvironmentEnabled(ctx, projectId, userId)
	if err != nil {
		return false, err
	}
	return enabled, nil
}
