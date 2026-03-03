package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	SubscribeToDeployImage(handler func(deployment *value.ImageDeployment)) error

	SubscribeToDeployDatabase(handler func(deployment *value.Deployment)) error
	PublishDatabaseDeployed(deployment *value.DatabaseDeployment) error

	SubscribeToDeleteDeployment(handler func(deployment *value.Deployment)) error
	PublishDeploymentDeleted(deployment *value.DeploymentDeleted) error

	SubscribeToDeployIngress(handler func(deployment *value.IngressDeployment)) error
}
