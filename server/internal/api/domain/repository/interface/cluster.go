package interfaces

import (
	"context"
	"starliner.app/internal/api/domain/entity"
)

type ClusterRepository interface {
	CreateCluster(
		ctx context.Context,
		name string,
		serverType string,
		organizationId int64,
	) (*entity.Cluster, error)

	GetUserCluster(
		ctx context.Context,
		userId int64,
		clusterId int64,
	) (*entity.Cluster, error)

	GetCluster(
		ctx context.Context,
		clusterId int64,
	) (*entity.Cluster, error)

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
		status entity.ClusterStatus,
	) error

	UpdateClusterKubeconfig(
		ctx context.Context,
		id int64,
		kubeconfig *string,
	) error

	UpdateClusterLogs(
		ctx context.Context,
		id int64,
		logs string,
	) error

	GetUserClusterProvisioningLogs(
		ctx context.Context,
		userId int64,
		clusterId int64,
	) (*string, error)
}
