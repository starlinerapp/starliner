package pubsub

import (
	"context"
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

const BuildLogs natscore.Subject = "build.logs"

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

func (p *Pubsub) SubscribeToBuildLogs(
	ctx context.Context,
	buildId int64,
) (<-chan *coreValue.BuildLogChunk, func(), error) {
	ch := make(chan *coreValue.BuildLogChunk, 256)

	cancelSub, err := p.subscriber.SubscribeContext(
		ctx,
		BuildLogs,
		strconv.FormatInt(buildId, 10),
		func(msg []byte) {
			var chunk coreValue.BuildLogChunk
			if err := json.Unmarshal(msg, &chunk); err != nil {
				log.Printf("failed to unmarshal build log chunk: %v", err)
				return
			}
			select {
			case ch <- &chunk:
			case <-ctx.Done():
			}
		},
	)
	if err != nil {
		close(ch)
		return nil, func() {}, err
	}

	cancel := func() {
		cancelSub()
	}
	return ch, cancel, nil
}
