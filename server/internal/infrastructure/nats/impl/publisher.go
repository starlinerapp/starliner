package nats

import (
	"starliner.app/internal/domain/entity"
	"starliner.app/internal/infrastructure/nats"
	v1 "starliner.app/internal/infrastructure/nats/proto/v1"
	"strconv"
)

func (q *Queue) PublishBuildTriggered(build *entity.Build) error {
	return q.buildPublisher.Publish(nats.BuildTriggered, build.Id, &v1.Build{
		Id:             build.Id,
		Organization:   build.Organization,
		Project:        build.Project,
		Service:        build.Service,
		S3Key:          build.S3Key,
		RootDirectory:  build.RootDirectory,
		DockerfilePath: build.DockerfilePath,
	})
}

func (q *Queue) PublishCreateCluster(cluster *entity.Cluster) error {
	return q.clusterPublisher.Publish(nats.CreateCluster, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id:             cluster.Id,
		Name:           cluster.Name,
		OrganizationId: cluster.OrganizationId,
	})
}

func (q *Queue) PublishDeleteCluster(cluster *entity.Cluster) error {
	return q.clusterPublisher.Publish(nats.DeleteCluster, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id:             cluster.Id,
		Name:           cluster.Name,
		OrganizationId: cluster.OrganizationId,
	})
}

func (q *Queue) PublishCreateProject(project *entity.Project) error {
	return q.projectPublisher.Publish(nats.CreateProject, strconv.FormatInt(project.Id, 10), &v1.Project{
		Id:             project.Id,
		Name:           project.Name,
		OrganizationId: project.OrganizationId,
		ClusterId:      *project.ClusterId,
	})
}
