package port

import (
	"starliner.app/internal/domain/value"
)

type Queue interface {
	PublishBuildTriggered(build *value.BuildMessage) error
	PublishCreateCluster(cluster *value.ClusterMessage) error
	PublishDeleteCluster(cluster *value.ClusterMessage) error
	PublishDeployDatabase(deployment *value.DeploymentMessage) error

	SubscribeToBuildTriggered(handler func(build *value.BuildMessage)) error
	SubscribeToCreateCluster(handler func(cluster *value.ClusterMessage)) error
	SubscribeToDeleteCluster(handler func(cluster *value.ClusterMessage)) error
	SubscribeToDeployDatabase(handler func(deployment *value.DeploymentMessage)) error
}
