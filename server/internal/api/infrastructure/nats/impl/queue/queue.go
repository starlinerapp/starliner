package queue

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
	"strconv"
)

const (
	BuildTriggered    jetstream.Subject = "build.triggered"
	CreateCluster     jetstream.Subject = "create.cluster"
	ClusterCreated    jetstream.Subject = "cluster.created"
	DeleteCluster     jetstream.Subject = "delete.cluster"
	ClusterDeleted    jetstream.Subject = "cluster.deleted"
	DeployApplication jetstream.Subject = "deploy.application"
	DeployDatabase    jetstream.Subject = "deploy.database"
	DeleteDatabase    jetstream.Subject = "delete.database"
	DatabaseDeleted   jetstream.Subject = "database.deleted"
	DeployIngress     jetstream.Subject = "deploy.ingress"
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

func (q *Queue) PublishBuildTriggered(build *value.TriggerBuild) error {
	data, err := json.Marshal(build)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(BuildTriggered, build.Id, data)
}

func (q *Queue) PublishCreateCluster(cluster *value.ProvisionCluster) error {
	data, err := json.Marshal(cluster)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(CreateCluster, strconv.FormatInt(cluster.Id, 10), data)
}

func (q *Queue) SubscribeToClusterCreated(handler func(cluster *value.ClusterCreated)) error {
	return q.subscriber.Subscribe(ClusterCreated, "*", "clusterCreated", func(msg []byte) {
		var c value.ClusterCreated
		if err := json.Unmarshal(msg, &c); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&c)
	})
}

func (q *Queue) SubscribeToClusterDeleted(handler func(cluster *value.ClusterDeleted)) error {
	return q.subscriber.Subscribe(ClusterDeleted, "*", "clusterDeleted", func(msg []byte) {
		var c value.ClusterDeleted
		if err := json.Unmarshal(msg, &c); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&c)
	})
}

func (q *Queue) PublishDeleteCluster(cluster *value.DeleteCluster) error {
	data, err := json.Marshal(cluster)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeleteCluster, strconv.FormatInt(cluster.Id, 10), data)
}

func (q *Queue) PublishDeployApplication(deployment *value.ApplicationDeployment) error {
	data, err := json.Marshal(deployment)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeployApplication, "*", data)
}

func (q *Queue) PublishDeployDatabase(deployment *value.DatabaseDeployment) error {
	data, err := json.Marshal(deployment)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeployDatabase, "*", data)
}

func (q *Queue) PublishDeleteDatabase(deployment *value.DatabaseDeployment) error {
	data, err := json.Marshal(deployment)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeleteDatabase, "*", data)
}

func (q *Queue) SubscribeToDatabaseDeleted(handler func(deployment *value.DeploymentDeleted)) error {
	return q.subscriber.Subscribe(DatabaseDeleted, "*", "databaseDeleted", func(msg []byte) {
		var d value.DeploymentDeleted
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&d)
	})
}

func (q *Queue) PublishDeployIngress(deployment *value.IngressDeployment) error {
	data, err := json.Marshal(deployment)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeployIngress, "*", data)
}
