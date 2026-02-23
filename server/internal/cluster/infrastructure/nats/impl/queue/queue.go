package queue

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
	"strconv"
)

const (
	DeployImage     jetstream.Subject = "deploy.image"
	DeployDatabase  jetstream.Subject = "deploy.database"
	DeleteDatabase  jetstream.Subject = "delete.database"
	DatabaseDeleted jetstream.Subject = "database.deleted"
	DeployIngress   jetstream.Subject = "deploy.ingress"
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

func (q *Queue) SubscribeToDeployImage(handler func(deployment *value.ImageDeployment)) error {
	return q.subscriber.Subscribe(DeployImage, "*", "deployImage", func(msg []byte) {
		var d value.ImageDeployment
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&d)
	})
}

func (q *Queue) SubscribeToDeployDatabase(handler func(deployment *value.DatabaseDeployment)) error {
	return q.subscriber.Subscribe(DeployDatabase, "*", "deployDatabase", func(msg []byte) {
		var d value.DatabaseDeployment
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&d)
	})
}

func (q *Queue) SubscribeToDeleteDatabase(handler func(deployment *value.DatabaseDeployment)) error {
	return q.subscriber.Subscribe(DeleteDatabase, "*", "deleteDatabase", func(msg []byte) {
		var d value.DatabaseDeployment
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&d)
	})
}

func (q *Queue) SubscribeToDeployIngress(handler func(deployment *value.IngressDeployment)) error {
	return q.subscriber.Subscribe(DeployIngress, "*", "deployIngress", func(msg []byte) {
		var d value.IngressDeployment
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&d)
	})
}

func (q *Queue) PublishDatabaseDeleted(deployment *value.DeploymentDeleted) error {
	data, err := json.Marshal(deployment)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DatabaseDeleted, strconv.FormatInt(deployment.DeploymentId, 10), data)
}
