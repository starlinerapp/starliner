package service

import (
	"context"
	"errors"
	interfaces "starliner.app/internal/api/domain/repository/interface"
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
