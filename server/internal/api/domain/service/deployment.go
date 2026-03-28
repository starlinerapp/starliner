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
