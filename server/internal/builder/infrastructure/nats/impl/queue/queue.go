package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nats-io/nats.go"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
)

const (
	BuildTriggered jetstream.Subject = "build.triggered"
	BuildCompleted jetstream.Subject = "build.completed"
)

type Queue struct {
	subscriber *jetstream.Subscriber
	publisher  *jetstream.Publisher
}

func NewQueue(js nats.JetStreamContext) port.Queue {
	return &Queue{
		subscriber: jetstream.NewSubscriber(js),
		publisher:  jetstream.NewPublisher(js),
	}
}

func (q *Queue) SubscribeToBuildTriggered(handler func(build *value.TriggerBuild)) error {
	return q.subscriber.Subscribe(BuildTriggered, "*", "buildTriggered", func(msg []byte) {
		var b value.TriggerBuild
		if err := json.Unmarshal(msg, &b); err != nil {
			log.Printf("failed to unmarshal: %v", err)
		}
		hostname, _ := os.Hostname()
		log.Printf("[builder replica=%s] picked up build job: buildID=%d deploymentID=%d", hostname, b.BuildId, b.DeploymentId)
		handler(&b)
	})
}

func (q *Queue) PublishBuildCompleted(build *value.BuildCompleted) error {
	data, err := json.Marshal(build)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	return q.publisher.Publish(BuildCompleted, "*", data)
}
