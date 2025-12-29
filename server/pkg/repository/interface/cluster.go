package interfaces

import (
	"context"
	"starliner.app/pkg/domain"
)

type ClusterRepository interface {
	CreateCluster(
		ctx context.Context,
		name string,
		organizationId int64,
		ipv4Address *string,
		publicKey *string,
		privateKeyRef *string,
	) (*domain.Cluster, error)

	GetCluster(
		ctx context.Context,
		clusterId int64,
	) (*domain.Cluster, error)

	DeleteCluster(
		ctx context.Context,
		id int64,
	) error
}
