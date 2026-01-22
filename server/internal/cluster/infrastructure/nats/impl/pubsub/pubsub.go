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

type Pubsub struct {
	subscriber *natscore.Subscriber
}

func NewPubsub(conn *nats.Conn) port.Pubsub {
	return &Pubsub{
		subscriber: natscore.NewSubscriber(conn),
	}
}

func (p *Pubsub) SubscribeToDeploymentStatusRequest(handler func(deployment *value.Deployment)) error {
	return p.subscriber.Subscribe(DeploymentStatusRequest, "*", func(msg []byte) {
		var d value.Deployment
		if err := json.Unmarshal(msg, &d); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		handler(&d)
	})
}
