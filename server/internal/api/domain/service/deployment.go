package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	interfaces "starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/domain/value"
)

type DeploymentService struct {
	deploymentRepository interfaces.DeploymentRepository
}

func NewDeploymentService(deploymentRepository interfaces.DeploymentRepository) *DeploymentService {
	return &DeploymentService{
		deploymentRepository: deploymentRepository,
	}
}

func (ds *DeploymentService) ValidateUserPermission(ctx context.Context, userId int64, deploymentId int64) error {
	deployment, err := ds.deploymentRepository.GetUserDeployment(ctx, userId, deploymentId)
	if err != nil {
		return err
	}
	if deployment == nil {
		return errors.New("user not authorized")
	}

	return nil
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
