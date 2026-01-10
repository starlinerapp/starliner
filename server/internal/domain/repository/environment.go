package repository

import (
	"context"
	"starliner.app/internal/domain/entity"
	"starliner.app/internal/domain/repository/interface"
	"starliner.app/internal/infrastructure/postgres/sqlc"
	"starliner.app/internal/infrastructure/postgres/utils"
)

type EnvironmentRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.EnvironmentRepository = (*EnvironmentRepository)(nil)

func NewEnvironmentRepository(queries *sqlc.Queries) interfaces.EnvironmentRepository {
	return &EnvironmentRepository{queries: queries}
}

func (er *EnvironmentRepository) CreateEnvironment(ctx context.Context, name string, slug string, projectId int64) (*entity.Environment, error) {
	env, err := er.queries.CreateEnvironment(ctx, sqlc.CreateEnvironmentParams{
		Name:      name,
		Slug:      slug,
		ProjectID: projectId,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Environment{
		Id:   env.ID,
		Slug: env.Slug,
		Name: env.Name,
	}, nil
}

func (er *EnvironmentRepository) GetEnvironmentAuthorizedUsers(ctx context.Context, clusterId int64) (users []int64, err error) {
	users, err = er.queries.GetEnvironmentAuthorizedUsers(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (er *EnvironmentRepository) GetEnvironmentCluster(ctx context.Context, environmentId int64) (*entity.Cluster, error) {
	cluster, err := er.queries.GetEnvironmentCluster(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	return &entity.Cluster{
		Id:             cluster.ID,
		Name:           cluster.Name,
		Status:         entity.ClusterStatus(cluster.Status),
		IPv4Address:    utils.PtrFromNullString(cluster.Ipv4Address),
		PublicKey:      utils.PtrFromNullString(cluster.PublicKey),
		PrivateKey:     utils.PtrFromNullString(cluster.PrivateKey),
		PulumiStackId:  utils.PtrFromNullString(cluster.PulumiStackID),
		Kubeconfig:     utils.PtrFromNullString(cluster.Kubeconfig),
		OrganizationId: cluster.OrganizationID,
	}, nil
}
