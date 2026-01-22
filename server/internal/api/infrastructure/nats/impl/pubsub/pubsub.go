package pubsub

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"starliner.app/internal/api/domain/port"
	coreValue "starliner.app/internal/core/domain/value"
	natscore "starliner.app/internal/core/infrastructure/nats/core"
	"strconv"
)

const DeploymentStatusRequest natscore.Subject = "deployment.status.request"

type Pubsub struct {
	publisher *natscore.Publisher
}

func NewPubsub(conn *nats.Conn) port.Pubsub {
	return &Pubsub{
		publisher: natscore.NewPublisher(conn),
	}
}

func (p *Pubsub) PublishDeploymentStatusRequest(deployment coreValue.Deployment) error {
	d, err := json.Marshal(deployment)
	if err != nil {
		return err
	}
	return p.publisher.Publish(DeploymentStatusRequest, strconv.FormatInt(deployment.DeploymentId, 10), d)
}
