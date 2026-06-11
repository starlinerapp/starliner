package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/nats-io/nats.go"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
)

const (
	BuildTriggered                jetstream.Subject = "build.triggered"
	BuildSucceeded                jetstream.Subject = "build.succeeded"
	BuildFailed                   jetstream.Subject = "build.failed"
	CreateCluster                 jetstream.Subject = "create.cluster"
	DeleteCluster                 jetstream.Subject = "delete.cluster"
	ReconcileCluster              jetstream.Subject = "reconcile.cluster"
	DeployImage                   jetstream.Subject = "deploy.image"
	DeployDatabase                jetstream.Subject = "deploy.database"
	DeployIngress                 jetstream.Subject = "deploy.ingress"
	DeleteDeployment              jetstream.Subject = "delete.deployment"
	DeploymentStatusLogsCompleted jetstream.Subject = "deployment.status_logs.completed"

	ClusterProvisionedSuccess jetstream.Subject = "cluster.provisioned.success"
	ClusterProvisionedFailure jetstream.Subject = "cluster.provisioned.failure"
	ClusterDeletedSuccess     jetstream.Subject = "cluster.deleted.success"
	ClusterDeletedFailure     jetstream.Subject = "cluster.deleted.failure"

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

func (q *Queue) PublishBuildTriggered(build *value.TriggerBuild) error {
	data, err := json.Marshal(build)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(BuildTriggered, strconv.FormatInt(build.DeploymentId, 10), data)
}

func (q *Queue) SubscribeToBuildSucceeded(handler func(build *value.BuildSucceeded)) error {
	return q.subscriber.Subscribe(BuildSucceeded, "*", "buildSucceeded", func(msg []byte) {
		var b value.BuildSucceeded
		if err := json.Unmarshal(msg, &b); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&b)
	})
}

func (q *Queue) SubscribeToBuildFailed(handler func(build *value.BuildFailed)) error {
	return q.subscriber.Subscribe(BuildFailed, "*", "buildFailed", func(msg []byte) {
		var b value.BuildFailed
		if err := json.Unmarshal(msg, &b); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&b)
	})
}

func (q *Queue) PublishCreateCluster(cluster *value.ProvisionCluster) error {
	data, err := json.Marshal(cluster)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(CreateCluster, strconv.FormatInt(cluster.Id, 10), data)
}

func (q *Queue) SubscribeToClusterProvisionedSuccess(handler func(event *value.ClusterProvisionedSuccess)) error {
	return q.subscriber.Subscribe(ClusterProvisionedSuccess, "*", "clusterProvisionedSuccess", func(msg []byte) {
		var e value.ClusterProvisionedSuccess
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) SubscribeToClusterProvisionedFailure(handler func(event *value.ClusterProvisionedFailure)) error {
	return q.subscriber.Subscribe(ClusterProvisionedFailure, "*", "clusterProvisionedFailure", func(msg []byte) {
		var e value.ClusterProvisionedFailure
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) SubscribeToClusterDeletedSuccess(handler func(event *value.ClusterDeletedSuccess)) error {
	return q.subscriber.Subscribe(ClusterDeletedSuccess, "*", "clusterDeletedSuccess", func(msg []byte) {
		var e value.ClusterDeletedSuccess
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) SubscribeToClusterDeletedFailure(handler func(event *value.ClusterDeletedFailure)) error {
	return q.subscriber.Subscribe(ClusterDeletedFailure, "*", "clusterDeletedFailure", func(msg []byte) {
		var e value.ClusterDeletedFailure
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) PublishDeleteCluster(cluster *value.DeleteCluster) error {
	data, err := json.Marshal(cluster)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeleteCluster, strconv.FormatInt(cluster.Id, 10), data)
}

func (q *Queue) PublishReconcileCluster(cluster *value.ReconcileCluster) error {
	data, err := json.Marshal(cluster)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(ReconcileCluster, strconv.FormatInt(cluster.Id, 10), data)
}

func (q *Queue) PublishDeployImage(deployment *value.ImageDeployment) error {
	data, err := json.Marshal(deployment)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeployImage, "*", data)
}

func (q *Queue) SubscribeToImageDeployedSuccess(handler func(event *value.ImageDeployedSuccess)) error {
	return q.subscriber.Subscribe(ImageDeployedSuccess, "*", "imageDeployedSuccess", func(msg []byte) {
		var e value.ImageDeployedSuccess
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) SubscribeToImageDeployedFailure(handler func(event *value.ImageDeployedFailure)) error {
	return q.subscriber.Subscribe(ImageDeployedFailure, "*", "imageDeployedFailure", func(msg []byte) {
		var e value.ImageDeployedFailure
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) PublishDeployDatabase(deployment *value.Deployment) error {
	data, err := json.Marshal(deployment)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeployDatabase, "*", data)
}

func (q *Queue) SubscribeToDatabaseDeployedSuccess(handler func(event *value.DatabaseDeployedSuccess)) error {
	return q.subscriber.Subscribe(DatabaseDeployedSuccess, "*", "databaseDeployedSuccess", func(msg []byte) {
		var e value.DatabaseDeployedSuccess
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) SubscribeToDatabaseDeployedFailure(handler func(event *value.DatabaseDeployedFailure)) error {
	return q.subscriber.Subscribe(DatabaseDeployedFailure, "*", "databaseDeployedFailure", func(msg []byte) {
		var e value.DatabaseDeployedFailure
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) PublishDeleteDeployment(deployment *value.Deployment) error {
	data, err := json.Marshal(deployment)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeleteDeployment, "*", data)
}

func (q *Queue) SubscribeToDeploymentDeletedSuccess(handler func(event *value.DeploymentDeletedSuccess)) error {
	return q.subscriber.Subscribe(DeploymentDeletedSuccess, "*", "deploymentDeletedSuccess", func(msg []byte) {
		var e value.DeploymentDeletedSuccess
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) SubscribeToDeploymentDeletedFailure(handler func(event *value.DeploymentDeletedFailure)) error {
	return q.subscriber.Subscribe(DeploymentDeletedFailure, "*", "deploymentDeletedFailure", func(msg []byte) {
		var e value.DeploymentDeletedFailure
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) PublishDeployIngress(deployment *value.IngressDeployment) error {
	data, err := json.Marshal(deployment)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return q.publisher.Publish(DeployIngress, "*", data)
}

func (q *Queue) SubscribeToIngressDeployedSuccess(handler func(event *value.IngressDeployedSuccess)) error {
	return q.subscriber.Subscribe(IngressDeployedSuccess, "*", "ingressDeployedSuccess", func(msg []byte) {
		var e value.IngressDeployedSuccess
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) SubscribeToIngressDeployedFailure(handler func(event *value.IngressDeployedFailure)) error {
	return q.subscriber.Subscribe(IngressDeployedFailure, "*", "ingressDeployedFailure", func(msg []byte) {
		var e value.IngressDeployedFailure
		if err := json.Unmarshal(msg, &e); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&e)
	})
}

func (q *Queue) SubscribeToDeploymentStatusLogsCompleted(handler func(completed *value.DeploymentStatusLogsCompleted)) error {
	return q.subscriber.Subscribe(DeploymentStatusLogsCompleted, "*", "deploymentStatusLogsCompleted", func(msg []byte) {
		var completed value.DeploymentStatusLogsCompleted
		if err := json.Unmarshal(msg, &completed); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&completed)
	})
}

