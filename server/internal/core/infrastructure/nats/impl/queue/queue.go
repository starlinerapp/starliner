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
	ClusterCreated nats.Subject = "cluster.created"
	DeleteCluster  nats.Subject = "delete.cluster"
	ClusterDeleted nats.Subject = "cluster.deleted"
	DeployDatabase nats.Subject = "deploy.database"
)

type Queue struct {
	buildPublisher       *nats.Publisher[*v1.Build]
	buildSubscriber      *nats.Subscriber[*v1.Build]
	clusterPublisher     *nats.Publisher[*v1.Cluster]
	clusterSubscriber    *nats.Subscriber[*v1.Cluster]
	deploymentPublisher  *nats.Publisher[*v1.Deployment]
	deploymentSubscriber *nats.Subscriber[*v1.Deployment]
}

func NewQueue(js natsgo.JetStreamContext) port.Queue {
	return &Queue{
		buildPublisher:       nats.NewPublisher[*v1.Build](js),
		buildSubscriber:      nats.NewSubscriber[*v1.Build](js),
		clusterPublisher:     nats.NewPublisher[*v1.Cluster](js),
		clusterSubscriber:    nats.NewSubscriber[*v1.Cluster](js),
		deploymentPublisher:  nats.NewPublisher[*v1.Deployment](js),
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

func (q *Queue) PublishCreateCluster(cluster *value.ProvisionCluster) error {
	return q.clusterPublisher.Publish(CreateCluster, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id:               cluster.Id,
		Name:             cluster.Name,
		OrganizationName: cluster.OrganizationName,
	})
}

func (q *Queue) SubscribeToCreateCluster(handler func(cluster *value.ProvisionCluster)) error {
	return q.clusterSubscriber.Subscribe(CreateCluster, "*", "createCluster", func(cluster *v1.Cluster) {
		handler(&value.ProvisionCluster{
			Id:               cluster.Id,
			Name:             cluster.Name,
			OrganizationName: cluster.OrganizationName,
		})
	})
}

func (q *Queue) PublishClusterCreated(cluster *value.ClusterCreated) error {
	return q.clusterPublisher.Publish(ClusterCreated, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id:               cluster.Id,
		ProvisioningId:   cluster.ProvisioningId,
		Ipv4Address:      cluster.IPv4Address,
		PublicKey:        cluster.PublicKey,
		PrivateKey:       cluster.PrivateKey,
		KubeconfigBase64: cluster.KubeconfigBase64,
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

func (q *Queue) PublishClusterDeleted(cluster *value.ClusterDeleted) error {
	return q.clusterPublisher.Publish(ClusterDeleted, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id: cluster.Id,
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

func (q *Queue) SubscribeToDeleteCluster(handler func(cluster *value.DeleteCluster)) error {
	return q.clusterSubscriber.Subscribe(DeleteCluster, "*", "deleteCluster", func(cluster *v1.Cluster) {
		handler(&value.DeleteCluster{
			Id:             cluster.Id,
			ProvisioningId: cluster.ProvisioningId,
		})
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
		DeploymentId:     deployment.DeploymentId,
		KubeconfigBase64: deployment.KubeconfigBase64,
		Database:         protoDB,
	})
}

func (q *Queue) SubscribeToDeployDatabase(handler func(deployment *value.Deployment)) error {
	return q.deploymentSubscriber.Subscribe(DeployDatabase, "*", "deployDatabase", func(cluster *v1.Deployment) {
		handler(&value.Deployment{
			DeploymentId:     cluster.DeploymentId,
			KubeconfigBase64: cluster.KubeconfigBase64,
			Database:         value.Database(cluster.Database),
		})
	})
}
