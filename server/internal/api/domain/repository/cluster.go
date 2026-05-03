package repository

import (
	"context"
	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
	"starliner.app/internal/api/infrastructure/postgres/utils"
)

type ClusterRepository struct {
	queries *sqlc.Queries
}

var _ interfaces.ClusterRepository = (*ClusterRepository)(nil)

func NewClusterRepository(queries *sqlc.Queries) interfaces.ClusterRepository {
	return &ClusterRepository{queries: queries}
}

func (cr *ClusterRepository) GetCluster(ctx context.Context, clusterId int64) (*entity.Cluster, error) {
	c, err := cr.queries.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	return &entity.Cluster{
		Id:             c.ID,
		Name:           c.Name,
		Status:         entity.ClusterStatus(c.Status),
		IPv4Address:    utils.PtrFromNullString(c.Ipv4Address),
		PublicKey:      utils.PtrFromNullString(c.PublicKey),
		PrivateKey:     utils.PtrFromNullString(c.PrivateKey),
		ProvisioningId: utils.PtrFromNullString(c.ProvisioningID),
		Kubeconfig:     utils.PtrFromNullString(c.Kubeconfig),
		Logs:           utils.PtrFromNullString(c.Logs),
		OrganizationId: c.OrganizationID,
		ServerType:     entity.ServerType(c.ServerType),
	}, nil
}

func (cr *ClusterRepository) GetUserCluster(ctx context.Context, userId int64, clusterId int64) (*entity.Cluster, error) {
	cluster, err := cr.queries.GetUserCluster(ctx, sqlc.GetUserClusterParams{
		UserID: userId,
		ID:     clusterId,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Cluster{
		Id:             cluster.ID,
		Name:           cluster.Name,
		Status:         entity.ClusterStatus(cluster.Status),
		User:           cluster.User,
		IPv4Address:    utils.PtrFromNullString(cluster.Ipv4Address),
		PublicKey:      utils.PtrFromNullString(cluster.PublicKey),
		PrivateKey:     utils.PtrFromNullString(cluster.PrivateKey),
		ProvisioningId: utils.PtrFromNullString(cluster.ProvisioningID),
		OrganizationId: cluster.OrganizationID,
		ServerType:     entity.ServerType(cluster.ServerType),
	}, nil
}

func (cr *ClusterRepository) CreateCluster(
	ctx context.Context,
	name string,
	serverType string,
	organizationId int64,
) (*entity.Cluster, error) {
	cluster, err := cr.queries.CreateCluster(ctx, sqlc.CreateClusterParams{
		Name:           name,
		ServerType:     serverType,
		OrganizationID: organizationId,
	})
	if err != nil {
		return nil, err
	}

	return &entity.Cluster{
		Id:             cluster.ID,
		Name:           cluster.Name,
		ServerType:     entity.ServerType(cluster.ServerType),
		Status:         entity.ClusterStatus(cluster.Status),
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
	return cr.queries.UpdateClusterProvisioningId(ctx, sqlc.UpdateClusterProvisioningIdParams{
		ProvisioningID: utils.NullStringFromPtr(pulumiStackId),
		ID:             id,
	})
}

func (cr *ClusterRepository) DeleteCluster(
	ctx context.Context,
	id int64,
) error {
	return cr.queries.DeleteCluster(ctx, id)
}

func (cr *ClusterRepository) UpdateClusterStatus(
	ctx context.Context,
	id int64,
	status entity.ClusterStatus,
) error {
	return cr.queries.UpdateClusterStatus(ctx, sqlc.UpdateClusterStatusParams{
		Status: sqlc.ClusterStatus(status),
		ID:     id,
	})
}

func (cr *ClusterRepository) UpdateClusterKubeconfig(
	ctx context.Context,
	id int64,
	kubeconfig *string,
) error {
	return cr.queries.UpdateClusterKubeconfig(ctx, sqlc.UpdateClusterKubeconfigParams{
		Kubeconfig: utils.NullStringFromPtr(kubeconfig),
		ID:         id,
	})
}

func (cr *ClusterRepository) UpdateClusterLogs(
	ctx context.Context,
	id int64,
	logs string,
) error {
	return cr.queries.UpdateClusterLogs(ctx, sqlc.UpdateClusterLogsParams{
		Logs: utils.NullStringFromPtr(&logs),
		ID:   id,
	})
}

func (cr *ClusterRepository) GetUserClusterProvisioningLogs(
	ctx context.Context,
	userId int64,
	clusterId int64,
) (*string, error) {
	res, err := cr.queries.GetUserClusterProvisioningLogs(ctx, sqlc.GetUserClusterProvisioningLogsParams{
		ClusterID: clusterId,
		UserID:    userId,
	})
	return utils.PtrFromNullString(res), err
}
