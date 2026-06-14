package pubsub

import (
	"encoding/json"
	"strconv"

	"github.com/nats-io/nats.go"
	"starliner.app/internal/cluster/domain/port"
	"starliner.app/internal/core/domain/value"
	natscore "starliner.app/internal/core/infrastructure/nats/core"
)

const ReconcileClusterRequest natscore.Subject = "reconcile.cluster.request"

type Pubsub struct {
	publisher *natscore.Publisher
}

func NewPubsub(conn *nats.Conn) port.Pubsub {
	return &Pubsub{
		publisher: natscore.NewPublisher(conn),
	}
}

func (p *Pubsub) PublishReconcileClusterRequest(request *value.ReconcileClusterRequest) error {
	d, err := json.Marshal(request)
	if err != nil {
		return err
	}
	return p.publisher.Publish(ReconcileClusterRequest, strconv.FormatInt(request.ClusterId, 10), d)
}
