package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	SubscribeToCreateCluster(handler func(cluster *value.ProvisionCluster)) error
	SubscribeToDeleteCluster(handler func(cluster *value.DeleteCluster)) error

	PublishClusterCreated(cluster *value.ClusterCreated) error
	PublishClusterDeleted(cluster *value.ClusterDeleted) error
}
