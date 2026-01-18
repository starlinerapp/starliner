package queue

import (
	natsgo "github.com/nats-io/nats.go"
	"starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats"
	"starliner.app/internal/core/infrastructure/nats/proto/v1"
	"strconv"
)

const (
	BuildTriggered nats.Subject = "build.triggered"
	CreateCluster  nats.Subject = "create.cluster"
	DeleteCluster  nats.Subject = "delete.cluster"
	DeployDatabase nats.Subject = "deploy.database"
)

type Queue struct {
	buildPublisher       *nats.Publisher[*v1.Build]
	clusterPublisher     *nats.Publisher[*v1.Cluster]
	deploymentPublisher  *nats.Publisher[*v1.Deployment]
	buildSubscriber      *nats.Subscriber[*v1.Build]
	clusterSubscriber    *nats.Subscriber[*v1.Cluster]
	deploymentSubscriber *nats.Subscriber[*v1.Deployment]
}

func NewQueue(js natsgo.JetStreamContext) port.Queue {
	return &Queue{
		buildPublisher:       nats.NewPublisher[*v1.Build](js),
		clusterPublisher:     nats.NewPublisher[*v1.Cluster](js),
		deploymentPublisher:  nats.NewPublisher[*v1.Deployment](js),
		buildSubscriber:      nats.NewSubscriber[*v1.Build](js),
		clusterSubscriber:    nats.NewSubscriber[*v1.Cluster](js),
		deploymentSubscriber: nats.NewSubscriber[*v1.Deployment](js),
	}
}

func (q *Queue) PublishBuildTriggered(build *value.Build) error {
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

func (q *Queue) PublishCreateCluster(cluster *value.Cluster) error {
	return q.clusterPublisher.Publish(CreateCluster, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id:               cluster.Id,
		Name:             cluster.Name,
		OrganizationName: cluster.OrganizationName,
	})
}

func (q *Queue) PublishDeleteCluster(cluster *value.Cluster) error {
	return q.clusterPublisher.Publish(DeleteCluster, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id:               cluster.Id,
		Name:             cluster.Name,
		OrganizationName: cluster.OrganizationName,
	})
}

func (q *Queue) PublishDeployDatabase(deployment *value.Deployment) error {
	var valueToProto = map[value.Database]v1.DatabaseType{
		value.Postgres: v1.DatabaseType_POSTGRES,
	}
	protoDB := func() v1.DatabaseType {
		if db, ok := valueToProto[deployment.Database]; ok {
			return db
		}
		return v1.DatabaseType_UNSPECIFIED
	}()

	return q.deploymentPublisher.Publish(DeployDatabase, "*", &v1.Deployment{
		DeploymentId: deployment.DeploymentId,
		ClusterId:    deployment.ClusterId,
		Database:     protoDB,
	})
}

func (q *Queue) SubscribeToBuildTriggered(handler func(build *value.Build)) error {
	return q.buildSubscriber.Subscribe(BuildTriggered, "*", "buildTriggered", func(build *v1.Build) {
		handler(&value.Build{
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

func (q *Queue) SubscribeToCreateCluster(handler func(cluster *value.Cluster)) error {
	return q.clusterSubscriber.Subscribe(CreateCluster, "*", "createCluster", func(cluster *v1.Cluster) {
		handler(&value.Cluster{
			Id:               cluster.Id,
			Name:             cluster.Name,
			OrganizationName: cluster.OrganizationName,
		})
	})
}

func (q *Queue) SubscribeToDeleteCluster(handler func(cluster *value.Cluster)) error {
	return q.clusterSubscriber.Subscribe(DeleteCluster, "*", "deleteCluster", func(cluster *v1.Cluster) {
		handler(&value.Cluster{
			Id:               cluster.Id,
			Name:             cluster.Name,
			OrganizationName: cluster.OrganizationName,
		})
	})
}

func (q *Queue) SubscribeToDeployDatabase(handler func(deployment *value.Deployment)) error {
	return q.deploymentSubscriber.Subscribe(DeployDatabase, "*", "deployDatabase", func(cluster *v1.Deployment) {
		handler(&value.Deployment{
			DeploymentId: cluster.DeploymentId,
			ClusterId:    cluster.ClusterId,
			Database:     value.Database(cluster.Database),
		})
	})
}
