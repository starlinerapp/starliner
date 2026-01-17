package queue

import (
	natsgo "github.com/nats-io/nats.go"
	"starliner.app/internal/domain/entity"
	"starliner.app/internal/domain/port"
	"starliner.app/internal/infrastructure/nats"
	v1 "starliner.app/internal/infrastructure/nats/proto/v1"
	"strconv"
)

const (
	BuildTriggered nats.Subject = "build.triggered"
	CreateCluster  nats.Subject = "create.cluster"
	DeleteCluster  nats.Subject = "delete.cluster"
	CreateProject  nats.Subject = "create.project"
)

type Queue struct {
	buildPublisher    *nats.Publisher[*v1.Build]
	clusterPublisher  *nats.Publisher[*v1.Cluster]
	projectPublisher  *nats.Publisher[*v1.Project]
	buildSubscriber   *nats.Subscriber[*v1.Build]
	clusterSubscriber *nats.Subscriber[*v1.Cluster]
	projectSubscriber *nats.Subscriber[*v1.Project]
}

func NewQueue(js natsgo.JetStreamContext) port.Queue {
	return &Queue{
		buildPublisher:    nats.NewPublisher[*v1.Build](js),
		clusterPublisher:  nats.NewPublisher[*v1.Cluster](js),
		projectPublisher:  nats.NewPublisher[*v1.Project](js),
		buildSubscriber:   nats.NewSubscriber[*v1.Build](js),
		clusterSubscriber: nats.NewSubscriber[*v1.Cluster](js),
		projectSubscriber: nats.NewSubscriber[*v1.Project](js),
	}
}

func (q *Queue) PublishBuildTriggered(build *entity.Build) error {
	return q.buildPublisher.Publish(BuildTriggered, build.Id, &v1.Build{
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
	return q.clusterPublisher.Publish(CreateCluster, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id:             cluster.Id,
		Name:           cluster.Name,
		OrganizationId: cluster.OrganizationId,
	})
}

func (q *Queue) PublishDeleteCluster(cluster *entity.Cluster) error {
	return q.clusterPublisher.Publish(DeleteCluster, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id:             cluster.Id,
		Name:           cluster.Name,
		OrganizationId: cluster.OrganizationId,
	})
}

func (q *Queue) PublishCreateProject(project *entity.ProjectWithEnvironments) error {
	return q.projectPublisher.Publish(CreateProject, strconv.FormatInt(project.Id, 10), &v1.Project{
		Id:             project.Id,
		Name:           project.Name,
		OrganizationId: project.OrganizationId,
		ClusterId:      *project.ClusterId,
	})
}

func (q *Queue) SubscribeToBuildTriggered(handler func(build *entity.Build)) error {
	return q.buildSubscriber.Subscribe(BuildTriggered, "*", "buildTriggered", func(build *v1.Build) {
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
	return q.clusterSubscriber.Subscribe(CreateCluster, "*", "createCluster", func(cluster *v1.Cluster) {
		handler(&entity.Cluster{
			Id:             cluster.Id,
			Name:           cluster.Name,
			OrganizationId: cluster.OrganizationId,
		})
	})
}

func (q *Queue) SubscribeToDeleteCluster(handler func(cluster *entity.Cluster)) error {
	return q.clusterSubscriber.Subscribe(DeleteCluster, "*", "deleteCluster", func(cluster *v1.Cluster) {
		handler(&entity.Cluster{
			Id:             cluster.Id,
			Name:           cluster.Name,
			OrganizationId: cluster.OrganizationId,
		})
	})
}

func (q *Queue) SubscribeToCreateProject(handler func(project *entity.ProjectWithEnvironments)) error {
	return q.projectSubscriber.Subscribe(CreateProject, "*", "createProject", func(project *v1.Project) {
		handler(&entity.ProjectWithEnvironments{
			Id:             project.Id,
			Name:           project.Name,
			OrganizationId: project.OrganizationId,
			ClusterId:      &project.ClusterId,
		})
	})
}
