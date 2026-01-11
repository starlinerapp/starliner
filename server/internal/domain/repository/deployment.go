package repository

import (
	"context"
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

func (dr *DeploymentRepository) CreateDeployment(ctx context.Context, name string, environmentId int64) error {
	return dr.queries.CreateDeployment(ctx, sqlc.CreateDeploymentParams{
		Name:          name,
		EnvironmentID: environmentId,
	})
}
