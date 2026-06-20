package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/nats-io/nats.go"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
)

const (
	BuildTriggered      jetstream.Subject = "build.triggered"
	BuildCompleted      jetstream.Subject = "build.completed"
	RunnerStatusChanged jetstream.Subject = "runner.status_changed"
	RunnerDeleted       jetstream.Subject = "runner.deleted"
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

func (q *Queue) PublishRunnerStatusChanged(status *value.RunnerStatusChanged) error {
	data, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	return q.publisher.Publish(RunnerStatusChanged, strconv.FormatInt(status.RunnerId, 10), data)
}

func (q *Queue) SubscribeToRunnerDeleted(handler func(runner *value.RunnerDeleted)) error {
	return q.subscriber.Subscribe(RunnerDeleted, "*", "runnerDeleted", func(msg []byte) {
		var runner value.RunnerDeleted
		if err := json.Unmarshal(msg, &runner); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
		}
		handler(&runner)
	})
}
