package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	SubscribeToDeployImage(handler func(deployment *value.ImageDeployment)) error
	PublishImageDeployedSuccess(event *value.ImageDeployedSuccess) error
	PublishImageDeployedFailure(event *value.ImageDeployedFailure) error

	SubscribeToDeployDatabase(handler func(deployment *value.Deployment)) error
	PublishDatabaseDeployedSuccess(event *value.DatabaseDeployedSuccess) error
	PublishDatabaseDeployedFailure(event *value.DatabaseDeployedFailure) error

	SubscribeToDeleteDeployment(handler func(deployment *value.Deployment)) error
	PublishDeploymentDeletedSuccess(event *value.DeploymentDeletedSuccess) error
	PublishDeploymentDeletedFailure(event *value.DeploymentDeletedFailure) error

	SubscribeToDeployIngress(handler func(deployment *value.IngressDeployment)) error
	PublishIngressDeployedSuccess(event *value.IngressDeployedSuccess) error
	PublishIngressDeployedFailure(event *value.IngressDeployedFailure) error

	PublishEnableIngressTLS(deployment *value.IngressDeployment) error
	SubscribeToEnableIngressTLS(handler func(deployment *value.IngressDeployment)) error

	PublishDeploymentStatusLogsCompleted(completed *value.DeploymentStatusLogsCompleted) error
}
