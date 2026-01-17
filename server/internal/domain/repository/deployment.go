package repository

import (
	"context"
	"starliner.app/internal/domain/entity"
	interfaces "starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/infrastructure/postgres/sqlc"
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
	username string,
	password string,
	environmentId int64,
) (deployment *entity.Deployment, err error) {
	d, err := dr.queries.CreateDatabaseDeployment(ctx, sqlc.CreateDatabaseDeploymentParams{
		Name:          name,
		Port:          port,
		Username:      username,
		Password:      password,
		EnvironmentID: environmentId,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Deployment{
		Id:            d.DeploymentID,
		Name:          d.Name,
		EnvironmentId: d.EnvironmentID,
	}, nil
}
