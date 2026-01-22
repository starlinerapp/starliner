package queue

import (
	natsgo "github.com/nats-io/nats.go"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
	"starliner.app/internal/core/infrastructure/nats/proto/v1"
)

const (
	DeployDatabase  jetstream.Subject = "deploy.database"
	DeleteDatabase  jetstream.Subject = "delete.database"
	DatabaseDeleted jetstream.Subject = "database.deleted"
)

type Queue struct {
	buildPublisher       *jetstream.Publisher[*v1.Build]
	buildSubscriber      *jetstream.Subscriber[*v1.Build]
	clusterPublisher     *jetstream.Publisher[*v1.Cluster]
	clusterSubscriber    *jetstream.Subscriber[*v1.Cluster]
	deploymentPublisher  *jetstream.Publisher[*v1.Deployment]
	deploymentSubscriber *jetstream.Subscriber[*v1.Deployment]
}

func NewQueue(js natsgo.JetStreamContext) port.Queue {
	return &Queue{
		buildPublisher:       jetstream.NewPublisher[*v1.Build](js),
		buildSubscriber:      jetstream.NewSubscriber[*v1.Build](js),
		clusterPublisher:     jetstream.NewPublisher[*v1.Cluster](js),
		clusterSubscriber:    jetstream.NewSubscriber[*v1.Cluster](js),
		deploymentPublisher:  jetstream.NewPublisher[*v1.Deployment](js),
		deploymentSubscriber: jetstream.NewSubscriber[*v1.Deployment](js),
	}
}

func (q *Queue) SubscribeToDeployDatabase(handler func(deployment *value.Deployment)) error {
	return q.deploymentSubscriber.Subscribe(DeployDatabase, "*", "deployDatabase", func(cluster *v1.Deployment) {
		handler(&value.Deployment{
			DeploymentId:     cluster.DeploymentId,
			KubeconfigBase64: cluster.KubeconfigBase64,
		})
	})
}

func (q *Queue) SubscribeToDeleteDatabase(handler func(deployment *value.Deployment)) error {
	return q.deploymentSubscriber.Subscribe(DeleteDatabase, "*", "deleteDatabase", func(cluster *v1.Deployment) {
		handler(&value.Deployment{
			DeploymentId:     cluster.DeploymentId,
			KubeconfigBase64: cluster.KubeconfigBase64,
		})
	})
}

func (q *Queue) PublishDatabaseDeleted(deployment *value.DeploymentDeleted) error {
	return q.deploymentPublisher.Publish(DatabaseDeleted, "*", &v1.Deployment{
		DeploymentId: deployment.DeploymentId,
	})
}
