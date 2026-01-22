package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	PublishBuildTriggered(build *value.TriggerBuild) error

	PublishCreateCluster(cluster *value.ProvisionCluster) error
	SubscribeToClusterCreated(handler func(cluster *value.ClusterCreated)) error

	PublishDeleteCluster(cluster *value.DeleteCluster) error
	SubscribeToClusterDeleted(handler func(cluster *value.ClusterDeleted)) error

	PublishDeployDatabase(deployment *value.Deployment) error

	PublishDeleteDatabase(deployment *value.Deployment) error
	SubscribeToDatabaseDeleted(handler func(deployment *value.DeploymentDeleted)) error
}
