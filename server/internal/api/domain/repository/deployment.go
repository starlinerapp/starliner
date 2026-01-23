package repository

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
	"starliner.app/internal/api/infrastructure/postgres/utils"
)

type DeploymentRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.DeploymentRepository = (*DeploymentRepository)(nil)

func NewDeploymentRepository(queries *sqlc.Queries) interfaces.DeploymentRepository {
	return &DeploymentRepository{queries: queries}
}

func (dr *DeploymentRepository) CreateDatabaseDeployment(
	ctx context.Context,
	name string,
	port string,
	status string,
	username string,
	password string,
	environmentId int64,
) (deployment *entity.DatabaseDeployment, err error) {
	d, err := dr.queries.CreateDatabaseDeployment(ctx, sqlc.CreateDatabaseDeploymentParams{
		Name:          name,
		Port:          port,
		Status:        utils.NullStringFromPtr(&status),
		Username:      username,
		Password:      password,
		EnvironmentID: environmentId,
	})
	if err != nil {
		return nil, err
	}

	return &entity.DatabaseDeployment{
		Id:            d.DeploymentID,
		Name:          d.Name,
		EnvironmentId: d.EnvironmentID,
	}, nil
}

func (dr *DeploymentRepository) GetUserDeployment(ctx context.Context, userId int64, deploymentId int64) (*entity.Deployment, error) {
	res, err := dr.queries.GetUserDeployment(ctx, sqlc.GetUserDeploymentParams{
		DeploymentID: deploymentId,
		UserID:       userId,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Deployment{
		Id:            res.ID,
		Name:          res.Name,
		Port:          res.Port,
		EnvironmentId: res.EnvironmentID,
	}, nil
}

func (dr *DeploymentRepository) GetDeploymentCluster(ctx context.Context, deploymentId int64) (*entity.Cluster, error) {
	res, err := dr.queries.GetDeploymentCluster(ctx, deploymentId)
	if err != nil {
		return nil, err
	}

	return &entity.Cluster{
		Id:             res.ID,
		Name:           res.Name,
		Status:         entity.ClusterStatus(res.Status),
		IPv4Address:    utils.PtrFromNullString(res.Ipv4Address),
		PublicKey:      utils.PtrFromNullString(res.PublicKey),
		PrivateKey:     utils.PtrFromNullString(res.PrivateKey),
		ProvisioningId: utils.PtrFromNullString(res.ProvisioningID),
		Kubeconfig:     utils.PtrFromNullString(res.Kubeconfig),
		OrganizationId: res.OrganizationID,
	}, nil
}

func (dr *DeploymentRepository) DeleteDeployment(ctx context.Context, deploymentId int64) error {
	return dr.queries.DeleteDeployment(ctx, deploymentId)
}

func (dr *DeploymentRepository) GetAllDeploymentsWithKubeconfig(ctx context.Context) ([]*entity.DeploymentWithKubeconfig, error) {
	rows, err := dr.queries.GetDeploymentsWithKubeconfig(ctx)
	if err != nil {
		return nil, err
	}

	deployments := make([]*entity.DeploymentWithKubeconfig, len(rows))
	for i, d := range rows {
		deployments[i] = &entity.DeploymentWithKubeconfig{
			Deployment: entity.Deployment{
				Id:            d.ID,
				Name:          d.Name,
				Port:          d.Port,
				EnvironmentId: d.EnvironmentID,
			},
			Kubeconfig: utils.PtrFromNullString(d.Kubeconfig),
		}
	}
	return deployments, nil
}

func (dr *DeploymentRepository) UpdateDeploymentStatus(ctx context.Context, deploymentId int64, status string) error {
	return dr.queries.UpdateDeploymentStatus(ctx, sqlc.UpdateDeploymentStatusParams{
		Status: utils.NullStringFromPtr(&status),
		ID:     deploymentId,
	})
}
