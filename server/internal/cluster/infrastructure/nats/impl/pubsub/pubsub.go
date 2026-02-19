package pubsub

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
	natscore "starliner.app/internal/core/infrastructure/nats/core"
)

const DeploymentStatusRequest natscore.Subject = "deployment.status.request"
const DeploymentStatusResponse natscore.Subject = "deployment.status.response"

type Pubsub struct {
	subscriber *natscore.Subscriber
	publisher  *natscore.Publisher
}

func NewPubsub(conn *nats.Conn) port.Pubsub {
	return &Pubsub{
		subscriber: natscore.NewSubscriber(conn),
		publisher:  natscore.NewPublisher(conn),
	}
}

func (p *Pubsub) SubscribeToDeploymentStatusRequest(handler func(deployment *value.DatabaseDeployment)) error {
	return p.subscriber.Subscribe(DeploymentStatusRequest, "*", func(msg []byte) {
		var d value.DatabaseDeployment
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&d)
	})
}

func (p *Pubsub) PublishDeploymentStatusResponse(health *value.HealthStatus) error {
	d, err := json.Marshal(health)
	if err != nil {
		return err
	}
	return p.publisher.Publish(DeploymentStatusResponse, "*", d)
}
