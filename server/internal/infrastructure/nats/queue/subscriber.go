package queue

import (
	"starliner.app/internal/domain/entity"
	"starliner.app/internal/infrastructure/nats"
	v1 "starliner.app/internal/infrastructure/nats/proto/v1"
)

func (q *Queue) SubscribeToBuildTriggered(handler func(build *entity.Build)) error {
	return q.buildSubscriber.Subscribe(nats.BuildTriggered, "*", "buildTriggered", func(build *v1.Build) {
		handler(&entity.Build{
			Id:             build.Id,
			Organization:   build.Organization,
			Project:        build.Project,
			Service:        build.Service,
			S3Key:          build.S3Key,
			RootDirectory:  build.RootDirectory,
			DockerfilePath: build.DockerfilePath,
		})
	})
}

func (q *Queue) SubscribeToCreateCluster(handler func(cluster *entity.Cluster)) error {
	return q.clusterSubscriber.Subscribe(nats.CreateCluster, "*", "createCluster", func(cluster *v1.Cluster) {
		handler(&entity.Cluster{
			Id:             cluster.Id,
			Name:           cluster.Name,
			OrganizationId: cluster.OrganizationId,
		})
	})
}

func (q *Queue) SubscribeToDeleteCluster(handler func(cluster *entity.Cluster)) error {
	return q.clusterSubscriber.Subscribe(nats.DeleteCluster, "*", "deleteCluster", func(cluster *v1.Cluster) {
		handler(&entity.Cluster{
			Id:             cluster.Id,
			Name:           cluster.Name,
			OrganizationId: cluster.OrganizationId,
		})
	})
}

func (q *Queue) SubscribeToCreateProject(handler func(project *entity.Project)) error {
	return q.projectSubscriber.Subscribe(nats.CreateProject, "*", "createProject", func(project *v1.Project) {
		handler(&entity.Project{
			Id:             project.Id,
			Name:           project.Name,
			OrganizationId: project.OrganizationId,
			ClusterId:      &project.ClusterId,
		})
	})
}
