package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	SubscribeToCreateCluster(handler func(cluster *value.ProvisionCluster)) error
	SubscribeToDeleteCluster(handler func(cluster *value.DeleteCluster)) error
	SubscribeToReconcileCluster(handler func(cluster *value.ReconcileCluster)) error

	PublishClusterProvisionedSuccess(event *value.ClusterProvisionedSuccess) error
	PublishClusterProvisionedFailure(event *value.ClusterProvisionedFailure) error

	PublishClusterDeletedSuccess(event *value.ClusterDeletedSuccess) error
	PublishClusterDeletedFailure(event *value.ClusterDeletedFailure) error
}
