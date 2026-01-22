package queue

import (
	natsgo "github.com/nats-io/nats.go"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
	"starliner.app/internal/core/infrastructure/nats/proto/v1"
	"strconv"
)

const (
	BuildTriggered  jetstream.Subject = "build.triggered"
	CreateCluster   jetstream.Subject = "create.cluster"
	ClusterCreated  jetstream.Subject = "cluster.created"
	DeleteCluster   jetstream.Subject = "delete.cluster"
	ClusterDeleted  jetstream.Subject = "cluster.deleted"
	DeployDatabase  jetstream.Subject = "deploy.database"
	DeleteDatabase  jetstream.Subject = "delete.database"
	DatabaseDeleted jetstream.Subject = "database.deleted"
)

type Queue struct {
	buildPublisher       *jetstream.Publisher[*v1.Build]
	clusterPublisher     *jetstream.Publisher[*v1.Cluster]
	clusterSubscriber    *jetstream.Subscriber[*v1.Cluster]
	deploymentPublisher  *jetstream.Publisher[*v1.Deployment]
	deploymentSubscriber *jetstream.Subscriber[*v1.Deployment]
}

func NewQueue(js natsgo.JetStreamContext) port.Queue {
	return &Queue{
		buildPublisher:       jetstream.NewPublisher[*v1.Build](js),
		clusterPublisher:     jetstream.NewPublisher[*v1.Cluster](js),
		clusterSubscriber:    jetstream.NewSubscriber[*v1.Cluster](js),
		deploymentPublisher:  jetstream.NewPublisher[*v1.Deployment](js),
		deploymentSubscriber: jetstream.NewSubscriber[*v1.Deployment](js),
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

func (q *Queue) PublishCreateCluster(cluster *value.ProvisionCluster) error {
	return q.clusterPublisher.Publish(CreateCluster, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id:               cluster.Id,
		Name:             cluster.Name,
		OrganizationName: cluster.OrganizationName,
	})
}

func (q *Queue) SubscribeToClusterCreated(handler func(cluster *value.ClusterCreated)) error {
	return q.clusterSubscriber.Subscribe(ClusterCreated, "*", "clusterCreated", func(cluster *v1.Cluster) {
		handler(&value.ClusterCreated{
			Id:               cluster.Id,
			ProvisioningId:   cluster.ProvisioningId,
			IPv4Address:      cluster.Ipv4Address,
			PublicKey:        cluster.PublicKey,
			PrivateKey:       cluster.PrivateKey,
			KubeconfigBase64: cluster.KubeconfigBase64,
		})
	})
}

func (q *Queue) SubscribeToClusterDeleted(handler func(cluster *value.ClusterDeleted)) error {
	return q.clusterSubscriber.Subscribe(ClusterDeleted, "*", "clusterDeleted", func(cluster *v1.Cluster) {
		handler(&value.ClusterDeleted{
			Id: cluster.Id,
		})
	})
}

func (q *Queue) PublishDeleteCluster(cluster *value.DeleteCluster) error {
	return q.clusterPublisher.Publish(DeleteCluster, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id:             cluster.Id,
		ProvisioningId: cluster.ProvisioningId,
	})
}

func (q *Queue) PublishDeployDatabase(deployment *value.Deployment) error {
	return q.deploymentPublisher.Publish(DeployDatabase, "*", &v1.Deployment{
		DeploymentId:     deployment.DeploymentId,
		KubeconfigBase64: deployment.KubeconfigBase64,
	})
}

func (q *Queue) PublishDeleteDatabase(deployment *value.Deployment) error {
	return q.deploymentPublisher.Publish(DeleteDatabase, "*", &v1.Deployment{
		DeploymentId:     deployment.DeploymentId,
		KubeconfigBase64: deployment.KubeconfigBase64,
	})
}

func (q *Queue) SubscribeToDatabaseDeleted(handler func(deployment *value.DeploymentDeleted)) error {
	return q.deploymentSubscriber.Subscribe(DatabaseDeleted, "*", "databaseDeleted", func(deployment *v1.Deployment) {
		handler(&value.DeploymentDeleted{
			DeploymentId: deployment.DeploymentId,
		})
	})
}
