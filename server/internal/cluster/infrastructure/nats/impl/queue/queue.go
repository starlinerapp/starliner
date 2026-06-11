package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/nats-io/nats.go"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
)

const (
	DeployImage                   jetstream.Subject = "deploy.image"
	DeployDatabase                jetstream.Subject = "deploy.database"
	DeployIngress                 jetstream.Subject = "deploy.ingress"
	EnableIngressTLS              jetstream.Subject = "enable.ingress.tls"
	DeleteDeployment              jetstream.Subject = "delete.deployment"
	DeploymentStatusLogsCompleted jetstream.Subject = "deployment.status_logs.completed"

	DatabaseDeployedSuccess  jetstream.Subject = "database.deployed.success"
	DatabaseDeployedFailure  jetstream.Subject = "database.deployed.failure"
	ImageDeployedSuccess     jetstream.Subject = "image.deployed.success"
	ImageDeployedFailure     jetstream.Subject = "image.deployed.failure"
	IngressDeployedSuccess   jetstream.Subject = "ingress.deployed.success"
	IngressDeployedFailure   jetstream.Subject = "ingress.deployed.failure"
	DeploymentDeletedSuccess jetstream.Subject = "deployment.deleted.success"
	DeploymentDeletedFailure jetstream.Subject = "deployment.deleted.failure"
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

func (q *Queue) SubscribeToDeployDatabase(handler func(deployment *value.Deployment)) error {
	return q.subscriber.Subscribe(DeployDatabase, "*", "deployDatabase", func(msg []byte) {
		var d value.Deployment
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&d)
	})
}

func (q *Queue) SubscribeToDeleteDeployment(handler func(deployment *value.Deployment)) error {
	return q.subscriber.Subscribe(DeleteDeployment, "*", "deleteDeployment", func(msg []byte) {
		var d value.Deployment
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
			return
		}
		handler(&d)
	})
}

func (q *Queue) PublishEnableIngressTLS(deployment *value.IngressDeployment) error {
	data, err := json.Marshal(deployment)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(EnableIngressTLS, strconv.FormatInt(deployment.DeploymentId, 10), data)
}

func (q *Queue) SubscribeToEnableIngressTLS(handler func(deployment *value.IngressDeployment)) error {
	return q.subscriber.Subscribe(EnableIngressTLS, "*", "enableIngressTLS", func(msg []byte) {
		var d value.IngressDeployment
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&d)
	})
}

func (q *Queue) PublishDeploymentStatusLogsCompleted(completed *value.DeploymentStatusLogsCompleted) error {
	data, err := json.Marshal(completed)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeploymentStatusLogsCompleted, strconv.FormatInt(completed.DeploymentId, 10), data)
}

func (q *Queue) PublishDatabaseDeployedSuccess(event *value.DatabaseDeployedSuccess) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DatabaseDeployedSuccess, strconv.FormatInt(event.DeploymentId, 10), data)
}

func (q *Queue) PublishDatabaseDeployedFailure(event *value.DatabaseDeployedFailure) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DatabaseDeployedFailure, strconv.FormatInt(event.DeploymentId, 10), data)
}

func (q *Queue) PublishImageDeployedSuccess(event *value.ImageDeployedSuccess) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(ImageDeployedSuccess, strconv.FormatInt(event.DeploymentId, 10), data)
}

func (q *Queue) PublishImageDeployedFailure(event *value.ImageDeployedFailure) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(ImageDeployedFailure, strconv.FormatInt(event.DeploymentId, 10), data)
}

func (q *Queue) PublishIngressDeployedSuccess(event *value.IngressDeployedSuccess) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(IngressDeployedSuccess, strconv.FormatInt(event.DeploymentId, 10), data)
}

func (q *Queue) PublishIngressDeployedFailure(event *value.IngressDeployedFailure) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(IngressDeployedFailure, strconv.FormatInt(event.DeploymentId, 10), data)
}

func (q *Queue) PublishDeploymentDeletedSuccess(event *value.DeploymentDeletedSuccess) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeploymentDeletedSuccess, strconv.FormatInt(event.DeploymentId, 10), data)
}

func (q *Queue) PublishDeploymentDeletedFailure(event *value.DeploymentDeletedFailure) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeploymentDeletedFailure, strconv.FormatInt(event.DeploymentId, 10), data)
}
