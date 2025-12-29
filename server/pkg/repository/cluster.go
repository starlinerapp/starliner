package repository

import (
	"context"
	"starliner.app/pkg/db/sqlc"
	"starliner.app/pkg/db/utils"
	"starliner.app/pkg/domain"
	interfaces "starliner.app/pkg/repository/interface"
)

type ClusterRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.ClusterRepository = (*ClusterRepository)(nil)

func NewClusterRepository(queries *sqlc.Queries) interfaces.ClusterRepository {
	return &ClusterRepository{queries: queries}
}

func (cr *ClusterRepository) GetCluster(ctx context.Context, clusterId int64) (*domain.Cluster, error) {
	cluster, err := cr.queries.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	return &domain.Cluster{
		Id:             cluster.ID,
		Name:           cluster.Name,
		IPv4Address:    utils.PtrFromNullString(cluster.Ipv4Address),
		PublicKey:      utils.PtrFromNullString(cluster.PublicKey),
		PrivateKeyRef:  utils.PtrFromNullString(cluster.PrivateKeyRef),
		OrganizationId: cluster.OrganizationID,
	}, nil

}

func (cr *ClusterRepository) CreateCluster(
	ctx context.Context,
	name string,
	organizationId int64,
	ipv4Address *string,
	publicKey *string,
	privateKeyRef *string,
) (*domain.Cluster, error) {
	cluster, err := cr.queries.CreateCluster(ctx, sqlc.CreateClusterParams{
		Name:           name,
		Ipv4Address:    utils.NullStringFromPtr(ipv4Address),
		PublicKey:      utils.NullStringFromPtr(publicKey),
		PrivateKeyRef:  utils.NullStringFromPtr(privateKeyRef),
		OrganizationID: organizationId,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Cluster{
		Id:             cluster.ID,
		Name:           cluster.Name,
		IPv4Address:    utils.PtrFromNullString(cluster.Ipv4Address),
		PublicKey:      utils.PtrFromNullString(cluster.PublicKey),
		PrivateKeyRef:  utils.PtrFromNullString(cluster.PrivateKeyRef),
		OrganizationId: cluster.OrganizationID,
	}, nil
}

func (cr *ClusterRepository) DeleteCluster(
	ctx context.Context,
	id int64,
) error {
	return cr.queries.DeleteCluster(ctx, id)
}
