package repository

import (
	"context"

	"starliner.app/internal/api/domain/entity"
	"starliner.app/internal/api/domain/repository/interface"
	"starliner.app/internal/api/infrastructure/postgres/mapper"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
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
		IPv4Address:    mapper.ToPtrFromNullString(c.Ipv4Address),
		PublicKey:      mapper.ToPtrFromNullString(c.PublicKey),
		PrivateKey:     mapper.ToPtrFromNullString(c.PrivateKey),
		ProvisioningId: mapper.ToPtrFromNullString(c.ProvisioningID),
		Kubeconfig:     mapper.ToPtrFromNullString(c.Kubeconfig),
		Logs:           mapper.ToPtrFromNullString(c.Logs),
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
		IPv4Address:    mapper.ToPtrFromNullString(cluster.Ipv4Address),
		PublicKey:      mapper.ToPtrFromNullString(cluster.PublicKey),
		PrivateKey:     mapper.ToPtrFromNullString(cluster.PrivateKey),
		ProvisioningId: mapper.ToPtrFromNullString(cluster.ProvisioningID),
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
		PublicKey:  mapper.ToNullStringFromPtr(publicKey),
		PrivateKey: mapper.ToNullStringFromPtr(privateKey),
		ID:         id,
	})
}

func (cr *ClusterRepository) UpdateClusterIPv4Address(
	ctx context.Context,
	id int64,
	ipv4Address *string,
) error {
	return cr.queries.UpdateClusterIPv4Address(ctx, sqlc.UpdateClusterIPv4AddressParams{
		Ipv4Address: mapper.ToNullStringFromPtr(ipv4Address),
		ID:          id,
	})
}

func (cr *ClusterRepository) UpdateClusterPulumiStackId(
	ctx context.Context,
	id int64,
	pulumiStackId *string,
) error {
	return cr.queries.UpdateClusterProvisioningId(ctx, sqlc.UpdateClusterProvisioningIdParams{
		ProvisioningID: mapper.ToNullStringFromPtr(pulumiStackId),
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
		Kubeconfig: mapper.ToNullStringFromPtr(kubeconfig),
		ID:         id,
	})
}

func (cr *ClusterRepository) UpdateClusterLogs(
	ctx context.Context,
	id int64,
	logs string,
) error {
	return cr.queries.UpdateClusterLogs(ctx, sqlc.UpdateClusterLogsParams{
		Logs: mapper.ToNullStringFromPtr(&logs),
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
	return mapper.ToPtrFromNullString(res), err
}
