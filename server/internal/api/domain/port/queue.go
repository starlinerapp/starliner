package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	PublishBuildTriggered(build *value.TriggerBuild) error

	PublishCreateCluster(cluster *value.ProvisionCluster) error
	SubscribeToClusterCreated(handler func(cluster *value.ClusterCreated)) error

	PublishDeleteCluster(cluster *value.DeleteCluster) error
	SubscribeToClusterDeleted(handler func(cluster *value.ClusterDeleted)) error

	PublishDeployImage(deployment *value.ImageDeployment) error

	PublishDeployDatabase(deployment *value.Deployment) error
	SubscribeToDatabaseDeploymentCreated(handler func(databaseDeployment *value.DatabaseDeployment)) error

	PublishDeleteDeployment(deployment *value.Deployment) error
	SubscribeToDeploymentDeleted(handler func(deployment *value.DeploymentDeleted)) error

	PublishDeployIngress(deployment *value.IngressDeployment) error
}
