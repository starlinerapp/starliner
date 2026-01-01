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
	c, err := cr.queries.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	return &domain.Cluster{
		Id:             c.ID,
		Name:           c.Name,
		IPv4Address:    utils.PtrFromNullString(c.Ipv4Address),
		PublicKey:      utils.PtrFromNullString(c.PublicKey),
		PrivateKey:     utils.PtrFromNullString(c.PrivateKey),
		PulumiStackId:  utils.PtrFromNullString(c.PulumiStackID),
		OrganizationId: c.OrganizationID,
	}, nil
}

func (cr *ClusterRepository) GetUserCluster(ctx context.Context, userId int64, clusterId int64) (*domain.Cluster, error) {
	cluster, err := cr.queries.GetUserCluster(ctx, sqlc.GetUserClusterParams{
		OwnerID: userId,
		ID:      clusterId,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Cluster{
		Id:             cluster.ID,
		Name:           cluster.Name,
		IPv4Address:    utils.PtrFromNullString(cluster.Ipv4Address),
		PublicKey:      utils.PtrFromNullString(cluster.PublicKey),
		PrivateKey:     utils.PtrFromNullString(cluster.PrivateKey),
		PulumiStackId:  utils.PtrFromNullString(cluster.PulumiStackID),
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
		PrivateKey:     utils.NullStringFromPtr(privateKeyRef),
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
		PrivateKey:     utils.PtrFromNullString(cluster.PrivateKey),
		PulumiStackId:  utils.PtrFromNullString(cluster.PulumiStackID),
		OrganizationId: cluster.OrganizationID,
	}, nil
}

func (cr *ClusterRepository) UpdateClusterPublicPrivateKey(ctx context.Context, id int64, publicKey *string, privateKey *string) error {
	return cr.queries.UpdateClusterPublicPrivateKeys(ctx, sqlc.UpdateClusterPublicPrivateKeysParams{
		PublicKey:  utils.NullStringFromPtr(publicKey),
		PrivateKey: utils.NullStringFromPtr(privateKey),
		ID:         id,
	})
}

func (cr *ClusterRepository) UpdateClusterIPv4Address(
	ctx context.Context,
	id int64,
	ipv4Address *string,
) error {
	return cr.queries.UpdateClusterIPv4Address(ctx, sqlc.UpdateClusterIPv4AddressParams{
		Ipv4Address: utils.NullStringFromPtr(ipv4Address),
		ID:          id,
	})
}

func (cr *ClusterRepository) UpdateClusterPulumiStackId(
	ctx context.Context,
	id int64,
	pulumiStackId *string,
) error {
	return cr.queries.UpdateClusterPulumiStackId(ctx, sqlc.UpdateClusterPulumiStackIdParams{
		PulumiStackID: utils.NullStringFromPtr(pulumiStackId),
		ID:            id,
	})
}

func (cr *ClusterRepository) DeleteCluster(
	ctx context.Context,
	id int64,
) error {
	return cr.queries.DeleteCluster(ctx, id)
}
