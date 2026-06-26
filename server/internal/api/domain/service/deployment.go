package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"starliner.app/internal/api/domain/entity"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
)

type DeploymentService struct {
	deploymentRepository interfaces.DeploymentRepository
	environmentService   *EnvironmentService
}

func NewDeploymentService(
	deploymentRepository interfaces.DeploymentRepository,
	environmentService *EnvironmentService,
) *DeploymentService {
	return &DeploymentService{
		deploymentRepository: deploymentRepository,
		environmentService:   environmentService,
	}
}

func (ds *DeploymentService) ValidateUserPermission(ctx context.Context, userId int64, deploymentId int64) error {
	_, err := ds.AuthorizeDeploymentAccess(ctx, userId, deploymentId)
	return err
}

func (ds *DeploymentService) AuthorizeDeploymentAccess(
	ctx context.Context,
	userId int64,
	deploymentId int64,
) (*entity.Deployment, error) {
	deployment, err := ds.deploymentRepository.GetDeploymentById(ctx, deploymentId)
	if err != nil {
		return nil, err
	}
	if deployment == nil || deployment.EnvironmentId == nil {
		return nil, errors.New("user not authorized")
	}

	if err := ds.environmentService.ValidateUserPermission(ctx, userId, *deployment.EnvironmentId); err != nil {
		return nil, err
	}

	return deployment, nil
}

func (ds *DeploymentService) AuthorizeGitDeploymentAccess(
	ctx context.Context,
	userId int64,
	deploymentId int64,
	environmentId int64,
) (*entity.GitDeployment, error) {
	if err := ds.environmentService.ValidateUserPermission(ctx, userId, environmentId); err != nil {
		return nil, err
	}

	deployment, err := ds.deploymentRepository.GetGitDeploymentById(ctx, deploymentId)
	if err != nil {
		return nil, err
	}
	if deployment == nil || deployment.EnvironmentId == nil || *deployment.EnvironmentId != environmentId {
		return nil, fmt.Errorf("git deployment not found")
	}

	return deployment, nil
}

func (ds *DeploymentService) AuthorizeImageDeploymentAccess(
	ctx context.Context,
	userId int64,
	deploymentId int64,
	environmentId int64,
) (*entity.ImageDeployment, error) {
	if err := ds.environmentService.ValidateUserPermission(ctx, userId, environmentId); err != nil {
		return nil, err
	}

	deployment, err := ds.deploymentRepository.GetImageDeploymentById(ctx, deploymentId)
	if err != nil {
		return nil, err
	}
	if deployment == nil || deployment.EnvironmentId == nil || *deployment.EnvironmentId != environmentId {
		return nil, fmt.Errorf("image deployment not found")
	}

	return deployment, nil
}

func (ds *DeploymentService) AuthorizeDatabaseDeploymentAccess(
	ctx context.Context,
	userId int64,
	deploymentId int64,
	environmentId int64,
) (*entity.DatabaseDeployment, error) {
	if err := ds.environmentService.ValidateUserPermission(ctx, userId, environmentId); err != nil {
		return nil, err
	}

	deployment, err := ds.deploymentRepository.GetDatabaseDeploymentById(ctx, deploymentId)
	if err != nil {
		return nil, err
	}
	if deployment == nil || deployment.EnvironmentId == nil || *deployment.EnvironmentId != environmentId {
		return nil, fmt.Errorf("database deployment not found")
	}

	return deployment, nil
}

func (ds *DeploymentService) ValidateIngressHostsAvailable(
	ctx context.Context,
	hosts []*value.IngressHost,
) error {

	var duplicates []string

	for _, h := range hosts {
		if h == nil {
			continue
		}

		found, err := ds.deploymentRepository.GetIngressHostByName(ctx, h.Host)
		if err != nil {
			return err
		}

		if found != nil {
			duplicates = append(duplicates, h.Host)
		}
	}

	if len(duplicates) > 0 {
		return fmt.Errorf("%w: %s",
			value.ErrIngressHostAlreadyExists,
			strings.Join(duplicates, ", "),
		)
	}

	return nil
}

func (ds *DeploymentService) getIngressHostSuffix(
	organizationSlug string,
	serverEnvironment string,
	deploymentDomain string,
) string {
	subdomain := ""
	switch serverEnvironment {
	case "local":
		subdomain = "dev"
	case "staging":
		subdomain = "staging"
	}

	if subdomain != "" {
		return "." + organizationSlug + "." + subdomain + "." + deploymentDomain
	}

	return "." + organizationSlug + "." + deploymentDomain
}

func (ds *DeploymentService) buildFullIngressHost(
	prefix value.IngressHostPrefix,
	organizationSlug string,
	serverEnvironment string,
	deploymentDomain string,
) string {
	return string(prefix) + ds.getIngressHostSuffix(
		organizationSlug,
		serverEnvironment,
		deploymentDomain,
	)
}

func (ds *DeploymentService) BuildIngressHosts(
	inputs []*value.IngressHostInput,
	organizationSlug string,
	serverEnvironment string,
	deploymentDomain string,
) ([]*value.IngressHost, error) {
	out := make([]*value.IngressHost, 0, len(inputs))

	for _, input := range inputs {
		if input == nil {
			continue
		}

		prefix, err := value.NewIngressHostPrefix(input.Prefix)
		if err != nil {
			return nil, err
		}

		out = append(out, &value.IngressHost{
			Host:  ds.buildFullIngressHost(prefix, organizationSlug, serverEnvironment, deploymentDomain),
			Paths: input.Paths,
		})
	}

	return out, nil
}
