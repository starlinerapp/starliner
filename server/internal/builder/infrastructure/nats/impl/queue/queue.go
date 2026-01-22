package queue

import (
	"github.com/nats-io/nats.go"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
	"starliner.app/internal/core/infrastructure/nats/jetstream"
	"starliner.app/internal/core/infrastructure/nats/proto/v1"
)

const (
	BuildTriggered jetstream.Subject = "build.triggered"
)

type Queue struct {
	buildSubscriber *jetstream.Subscriber[*v1.Build]
}

func NewQueue(js nats.JetStreamContext) port.Queue {
	return &Queue{
		buildSubscriber: jetstream.NewSubscriber[*v1.Build](js),
	}
}

func (q *Queue) SubscribeToBuildTriggered(handler func(build *value.Build)) error {
	return q.buildSubscriber.Subscribe(BuildTriggered, "*", "buildTriggered", func(build *v1.Build) {
		handler(&value.Build{
			Id:             build.Id,
			Organization:   build.Organization,
			Project:        build.Project,
			Service:        build.Service,
			S3Key:          build.S3Key,
			RootDirectory:  build.RootDirectory,
			DockerfilePath: build.DockerfilePath,
		})
	})
}
