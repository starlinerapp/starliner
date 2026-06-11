package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/nats-io/nats.go"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
	"starliner.app/internal/provisioner/domain/port"
)

const (
	CreateCluster   jetstream.Subject = "create.cluster"
	DeleteCluster   jetstream.Subject = "delete.cluster"
	ReconcileCluster jetstream.Subject = "reconcile.cluster"

	ClusterProvisionedSuccess jetstream.Subject = "cluster.provisioned.success"
	ClusterProvisionedFailure jetstream.Subject = "cluster.provisioned.failure"
	ClusterDeletedSuccess     jetstream.Subject = "cluster.deleted.success"
	ClusterDeletedFailure     jetstream.Subject = "cluster.deleted.failure"
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

func (q *Queue) SubscribeToDeleteCluster(handler func(cluster *value.DeleteCluster)) error {
	return q.subscriber.Subscribe(DeleteCluster, "*", "deleteCluster", func(cluster []byte) {
		var c value.DeleteCluster
		if err := json.Unmarshal(cluster, &c); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&c)
	})
}

func (q *Queue) SubscribeToReconcileCluster(handler func(cluster *value.ReconcileCluster)) error {
	return q.subscriber.Subscribe(ReconcileCluster, "*", "reconcileCluster", func(cluster []byte) {
		var c value.ReconcileCluster
		if err := json.Unmarshal(cluster, &c); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&c)
	})
}

func (q *Queue) PublishClusterProvisionedSuccess(event *value.ClusterProvisionedSuccess) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	return q.publisher.Publish(ClusterProvisionedSuccess, strconv.FormatInt(event.ClusterId, 10), data)
}

func (q *Queue) PublishClusterProvisionedFailure(event *value.ClusterProvisionedFailure) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	return q.publisher.Publish(ClusterProvisionedFailure, strconv.FormatInt(event.ClusterId, 10), data)
}

func (q *Queue) PublishClusterDeletedSuccess(event *value.ClusterDeletedSuccess) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	return q.publisher.Publish(ClusterDeletedSuccess, strconv.FormatInt(event.ClusterId, 10), data)
}

func (q *Queue) PublishClusterDeletedFailure(event *value.ClusterDeletedFailure) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	return q.publisher.Publish(ClusterDeletedFailure, strconv.FormatInt(event.ClusterId, 10), data)
}
