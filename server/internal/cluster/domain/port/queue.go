package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	SubscribeToDeployDatabase(handler func(deployment *value.Deployment)) error
	SubscribeToDeleteDatabase(handler func(deployment *value.Deployment)) error
	PublishDatabaseDeleted(deployment *value.DeploymentDeleted) error
}
