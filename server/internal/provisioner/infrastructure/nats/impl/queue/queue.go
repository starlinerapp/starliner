package queue

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
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
	publisher  *jetstream.Publisher
	subscriber *jetstream.Subscriber
}

func NewQueue(js nats.JetStreamContext) port.Queue {
	return &Queue{
		publisher:  jetstream.NewPublisher(js),
		subscriber: jetstream.NewSubscriber(js),
	}
}

func (q *Queue) SubscribeToCreateCluster(handler func(cluster *value.ProvisionCluster)) error {
	return q.subscriber.Subscribe(CreateCluster, "*", "createCluster", func(cluster []byte) {
		var c value.ProvisionCluster
		if err := json.Unmarshal(cluster, &c); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&c)
	})
}

func (q *Queue) PublishClusterCreated(cluster *value.ClusterCreated) error {
	data, err := json.Marshal(cluster)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	return q.publisher.Publish(ClusterCreated, strconv.FormatInt(cluster.Id, 10), data)
}

func (q *Queue) PublishClusterDeleted(cluster *value.ClusterDeleted) error {
	data, err := json.Marshal(cluster)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	return q.publisher.Publish(ClusterDeleted, strconv.FormatInt(cluster.Id, 10), data)
}

func (q *Queue) SubscribeToDeleteCluster(handler func(cluster *value.DeleteCluster)) error {
	return q.subscriber.Subscribe(DeleteCluster, "*", "deleteCluster", func(cluster []byte) {
		var c value.DeleteCluster
		if err := json.Unmarshal(cluster, &c); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&c)
	})
}
