package repository

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	sqlc2 "starliner.app/internal/core/infrastructure/postgres/sqlc"
)

type DeploymentRepository struct {
	queries *sqlc2.Queries
}

var _ _interface.DeploymentRepository = (*DeploymentRepository)(nil)

func NewDeploymentRepository(queries *sqlc2.Queries) _interface.DeploymentRepository {
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
	d, err := dr.queries.CreateDatabaseDeployment(ctx, sqlc2.CreateDatabaseDeploymentParams{
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
