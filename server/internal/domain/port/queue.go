package port

import "starliner.app/internal/domain/entity"

type Queue interface {
	PublishBuildTriggered(build *entity.Build) error
	PublishCreateCluster(cluster *entity.Cluster) error
	PublishDeleteCluster(cluster *entity.Cluster) error
	PublishCreateProject(project *entity.Project) error

	SubscribeToBuildTriggered(handler func(build *entity.Build)) error
	SubscribeToCreateCluster(handler func(cluster *entity.Cluster)) error
	SubscribeToDeleteCluster(handler func(cluster *entity.Cluster)) error
	SubscribeToCreateProject(handler func(project *entity.Project)) error
}
