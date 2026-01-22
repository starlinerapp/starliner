package queue

import (
	natsgo "github.com/nats-io/nats.go"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
	"starliner.app/internal/core/infrastructure/nats/proto/v1"
	"starliner.app/internal/provisioner/domain/port"
	"strconv"
)

const (
	CreateCluster  jetstream.Subject = "create.cluster"
	ClusterCreated jetstream.Subject = "cluster.created"
	DeleteCluster  jetstream.Subject = "delete.cluster"
	ClusterDeleted jetstream.Subject = "cluster.deleted"
)

type Queue struct {
	clusterPublisher  *jetstream.Publisher[*v1.Cluster]
	clusterSubscriber *jetstream.Subscriber[*v1.Cluster]
}

func NewQueue(js natsgo.JetStreamContext) port.Queue {
	return &Queue{
		clusterPublisher:  jetstream.NewPublisher[*v1.Cluster](js),
		clusterSubscriber: jetstream.NewSubscriber[*v1.Cluster](js),
	}
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

func (q *Queue) PublishClusterDeleted(cluster *value.ClusterDeleted) error {
	return q.clusterPublisher.Publish(ClusterDeleted, strconv.FormatInt(cluster.Id, 10), &v1.Cluster{
		Id: cluster.Id,
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
