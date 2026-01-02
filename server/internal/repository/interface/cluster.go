package interfaces

import (
	"context"
	"starliner.app/internal/domain"
)

type ClusterRepository interface {
	CreateCluster(
		ctx context.Context,
		name string,
		organizationId int64,
	) (*domain.Cluster, error)

	GetUserCluster(
		ctx context.Context,
		userId int64,
		clusterId int64,
	) (*domain.Cluster, error)

	GetCluster(
		ctx context.Context,
		clusterId int64,
	) (*domain.Cluster, error)

	DeleteCluster(
		ctx context.Context,
		id int64,
	) error

	UpdateClusterPublicPrivateKey(
		ctx context.Context,
		id int64,
		publicKey *string,
		privateKey *string,
	) error

	UpdateClusterIPv4Address(
		ctx context.Context,
		id int64,
		ipv4Address *string,
	) error

	UpdateClusterPulumiStackId(
		ctx context.Context,
		id int64,
		pulumiStackId *string,
	) error

	UpdateClusterStatus(
		ctx context.Context,
		id int64,
		status domain.ClusterStatus,
	) error
}
