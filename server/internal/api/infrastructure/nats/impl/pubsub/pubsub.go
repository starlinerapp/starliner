package pubsub

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"starliner.app/internal/api/domain/port"
	coreValue "starliner.app/internal/core/domain/value"
	natscore "starliner.app/internal/core/infrastructure/nats/core"
)

const ReconcileClusterRequest natscore.Subject = "reconcile.cluster.request"

type Pubsub struct {
	subscriber *natscore.Subscriber
}

func NewPubsub(conn *nats.Conn) port.Pubsub {
	return &Pubsub{
		subscriber: natscore.NewSubscriber(conn),
	}
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
