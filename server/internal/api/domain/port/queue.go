package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	PublishBuildTriggered(build *value.TriggerBuild) error
	SubscribeToBuildSucceeded(handler func(build *value.BuildSucceeded)) error
	SubscribeToBuildFailed(handler func(build *value.BuildFailed)) error

	PublishCreateCluster(cluster *value.ProvisionCluster) error
	SubscribeToClusterProvisionedSuccess(handler func(event *value.ClusterProvisionedSuccess)) error
	SubscribeToClusterProvisionedFailure(handler func(event *value.ClusterProvisionedFailure)) error

	PublishDeleteCluster(cluster *value.DeleteCluster) error
	PublishReconcileCluster(cluster *value.ReconcileCluster) error
	SubscribeToClusterDeletedSuccess(handler func(event *value.ClusterDeletedSuccess)) error
	SubscribeToClusterDeletedFailure(handler func(event *value.ClusterDeletedFailure)) error

	PublishDeployImage(deployment *value.ImageDeployment) error
	SubscribeToImageDeployedSuccess(handler func(event *value.ImageDeployedSuccess)) error
	SubscribeToImageDeployedFailure(handler func(event *value.ImageDeployedFailure)) error

	PublishDeployDatabase(deployment *value.Deployment) error
	SubscribeToDatabaseDeployedSuccess(handler func(event *value.DatabaseDeployedSuccess)) error
	SubscribeToDatabaseDeployedFailure(handler func(event *value.DatabaseDeployedFailure)) error

	PublishDeleteDeployment(deployment *value.Deployment) error
	SubscribeToDeploymentDeletedSuccess(handler func(event *value.DeploymentDeletedSuccess)) error
	SubscribeToDeploymentDeletedFailure(handler func(event *value.DeploymentDeletedFailure)) error

	PublishDeployIngress(deployment *value.IngressDeployment) error
	SubscribeToIngressDeployedSuccess(handler func(event *value.IngressDeployedSuccess)) error
	SubscribeToIngressDeployedFailure(handler func(event *value.IngressDeployedFailure)) error

	SubscribeToDeploymentStatusLogsCompleted(handler func(completed *value.DeploymentStatusLogsCompleted)) error
}
