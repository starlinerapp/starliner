package port

import (
	value2 "starliner.app/internal/core/domain/value"
)

type Queue interface {
	PublishBuildTriggered(build *value2.Build) error
	PublishCreateCluster(cluster *value2.Cluster) error
	PublishDeleteCluster(cluster *value2.Cluster) error
	PublishDeployDatabase(deployment *value2.Deployment) error

	SubscribeToBuildTriggered(handler func(build *value2.Build)) error
	SubscribeToCreateCluster(handler func(cluster *value2.Cluster)) error
	SubscribeToDeleteCluster(handler func(cluster *value2.Cluster)) error
	SubscribeToDeployDatabase(handler func(deployment *value2.Deployment)) error
}
