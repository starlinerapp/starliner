package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

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
	RunnerJob           jetstream.Subject = "runner.job"
	RunnerJobResult     jetstream.Subject = "runner.job.result"
)

const (
	RunnerJobsStream jetstream.Stream = "runner-jobs"
)

type Queue struct {
	js         nats.JetStreamContext
	subscriber *jetstream.Subscriber
	publisher  *jetstream.Publisher
	claimMu    sync.Mutex
	claimSubs  map[int64]*nats.Subscription
}

func NewQueue(js nats.JetStreamContext) port.Queue {
	return &Queue{
		js:         js,
		subscriber: jetstream.NewSubscriber(js),
		publisher:  jetstream.NewPublisher(js),
		claimSubs:  make(map[int64]*nats.Subscription),
	}
}

func (q *Queue) SubscribeToBuildTriggered(handler func(build *value.TriggerBuild)) error {
	return q.subscriber.Subscribe(BuildTriggered, "*", "buildTriggered", func(msg []byte) {
		var b value.TriggerBuild
		if err := json.Unmarshal(msg, &b); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			return
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

func (q *Queue) PublishRunnerJob(job *value.RunnerBuildJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal runner job: %w", err)
	}

	return q.publisher.Publish(RunnerJob, strconv.FormatInt(job.RunnerId, 10), data)
}

func (q *Queue) ClaimRunnerJob(runnerId int64) (*value.RunnerBuildJob, error) {
	q.claimMu.Lock()
	defer q.claimMu.Unlock()

	sub, err := q.ensureClaimSubscription(runnerId)
	if err != nil {
		return nil, err
	}

	msgs, err := sub.Fetch(1, nats.MaxWait(100*time.Millisecond))
	if err != nil {
		if err == nats.ErrTimeout {
			return nil, nil
		}
		return nil, fmt.Errorf("fetch runner job: %w", err)
	}
	if len(msgs) == 0 {
		return nil, nil
	}

	msg := msgs[0]
	var job value.RunnerBuildJob
	if err := json.Unmarshal(msg.Data, &job); err != nil {
		if termErr := msg.Term(); termErr != nil {
			log.Printf("failed to term invalid runner job message: %v", termErr)
		}
		return nil, fmt.Errorf("unmarshal runner job: %w", err)
	}

	if err := msg.Ack(); err != nil {
		if termErr := msg.Term(); termErr != nil {
			log.Printf("failed to term runner job after ack failure: %v", termErr)
		}
		return nil, fmt.Errorf("ack runner job: %w", err)
	}

	return &job, nil
}

func (q *Queue) ensureClaimSubscription(runnerId int64) (*nats.Subscription, error) {
	if sub, ok := q.claimSubs[runnerId]; ok && sub.IsValid() {
		return sub, nil
	}
	delete(q.claimSubs, runnerId)

	subject := fmt.Sprintf("%s.%d", RunnerJob, runnerId)
	consumerName := fmt.Sprintf("claim-%d", runnerId)
	stream := string(RunnerJobsStream)

	if _, err := q.js.ConsumerInfo(stream, consumerName); err == nats.ErrConsumerNotFound {
		_, err = q.js.AddConsumer(stream, &nats.ConsumerConfig{
			Durable:       consumerName,
			FilterSubject: subject,
			AckPolicy:     nats.AckExplicitPolicy,
			DeliverPolicy: nats.DeliverNewPolicy,
		})
		if err != nil {
			return nil, fmt.Errorf("create runner job consumer: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("consumer info runner job: %w", err)
	}

	sub, err := q.js.PullSubscribe(
		subject,
		consumerName,
		nats.BindStream(stream),
		nats.AckExplicit(),
	)
	if err != nil {
		return nil, fmt.Errorf("pull subscribe runner job: %w", err)
	}

	q.claimSubs[runnerId] = sub
	return sub, nil
}

func (q *Queue) PublishRunnerJobResult(result *value.RunnerBuildResult) error {
	payload, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal runner job result: %w", err)
	}

	return q.publisher.Publish(RunnerJobResult, strconv.FormatInt(result.BuildId, 10), payload)
}

func (q *Queue) SubscribeToRunnerJobResults(handler func(result *value.RunnerBuildResult)) error {
	return q.subscriber.Subscribe(RunnerJobResult, "*", "runnerJobResults", func(msg []byte) {
		var result value.RunnerBuildResult
		if err := json.Unmarshal(msg, &result); err != nil {
			log.Printf("failed to unmarshal runner job result: %v", err)
			return
		}
		handler(&result)
	})
}
