package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	SubscribeToDeployApplication(handler func(deployment *value.ApplicationDeployment)) error

	SubscribeToDeployDatabase(handler func(deployment *value.DatabaseDeployment)) error
	SubscribeToDeleteDatabase(handler func(deployment *value.DatabaseDeployment)) error
	PublishDatabaseDeleted(deployment *value.DeploymentDeleted) error

	SubscribeToDeployIngress(handler func(deployment *value.IngressDeployment)) error
}
