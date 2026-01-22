package port

import (
	"starliner.app/internal/core/domain/value"
)

type Queue interface {
	PublishBuildTriggered(build *value.Build) error
	SubscribeToBuildTriggered(handler func(build *value.Build)) error

	PublishCreateCluster(cluster *value.ProvisionCluster) error
	SubscribeToCreateCluster(handler func(cluster *value.ProvisionCluster)) error

	PublishClusterCreated(cluster *value.ClusterCreated) error
	SubscribeToClusterCreated(handler func(cluster *value.ClusterCreated)) error

	PublishClusterDeleted(cluster *value.ClusterDeleted) error
	SubscribeToClusterDeleted(handler func(cluster *value.ClusterDeleted)) error

	PublishDeleteCluster(cluster *value.DeleteCluster) error
	SubscribeToDeleteCluster(handler func(cluster *value.DeleteCluster)) error

	PublishDeployDatabase(deployment *value.Deployment) error
	SubscribeToDeployDatabase(handler func(deployment *value.Deployment)) error

	PublishDeleteDatabase(deployment *value.Deployment) error
	SubscribeToDeleteDatabase(handler func(deployment *value.Deployment)) error

	PublishDatabaseDeleted(deployment *value.DeploymentDeleted) error
	SubscribeToDatabaseDeleted(handler func(deployment *value.DeploymentDeleted)) error
}
