package pubsub

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/nats-io/nats.go"
	"starliner.app/internal/api/domain/port"
	coreValue "starliner.app/internal/core/domain/value"
	natscore "starliner.app/internal/core/infrastructure/nats/core"
)

const DeploymentStatusRequest natscore.Subject = "deployment.status.request"
const DeploymentStatusResponse natscore.Subject = "deployment.status.response"
const ReconcileClusterRequest natscore.Subject = "reconcile.cluster.request"

type Pubsub struct {
	publisher  *natscore.Publisher
	subscriber *natscore.Subscriber
}

func NewPubsub(conn *nats.Conn) port.Pubsub {
	return &Pubsub{
		publisher:  natscore.NewPublisher(conn),
		subscriber: natscore.NewSubscriber(conn),
	}
}

func (p *Pubsub) PublishDeploymentStatusRequest(deployment *coreValue.Deployment) error {
	d, err := json.Marshal(deployment)
	if err != nil {
		return err
	}
	return p.publisher.Publish(DeploymentStatusRequest, strconv.FormatInt(deployment.DeploymentId, 10), d)
}

func (p *Pubsub) SubscribeToDeploymentStatusResponse(handler func(health *coreValue.HealthStatus)) error {
	return p.subscriber.Subscribe(DeploymentStatusResponse, "*", func(msg []byte) {
		var health coreValue.HealthStatus
		if err := json.Unmarshal(msg, &health); err != nil {
			panic(err)
		}
		handler(&health)
	})
}

func (p *Pubsub) SubscribeToReconcileClusterRequest(handler func(request *coreValue.ReconcileClusterRequest)) error {
	return p.subscriber.Subscribe(ReconcileClusterRequest, "*", func(msg []byte) {
		var req coreValue.ReconcileClusterRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Printf("failed to unmarshal reconcile cluster request: %v", err)
			return
		}
		handler(&req)
	})
}
