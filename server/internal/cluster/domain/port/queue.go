package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	SubscribeToDeployDatabase(handler func(deployment *value.DatabaseDeployment)) error
	SubscribeToDeleteDatabase(handler func(deployment *value.DatabaseDeployment)) error
	PublishDatabaseDeleted(deployment *value.DeploymentDeleted) error

	SubscribeToDeployIngress(handler func(deployment *value.IngressDeployment)) error
}
