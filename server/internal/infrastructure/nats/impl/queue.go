package nats

import (
	natsgo "github.com/nats-io/nats.go"
	"starliner.app/internal/domain/port"
	"starliner.app/internal/infrastructure/nats"
	v1 "starliner.app/internal/infrastructure/nats/proto/v1"
)

type Queue struct {
	buildPublisher    *nats.Publisher[*v1.Build]
	clusterPublisher  *nats.Publisher[*v1.Cluster]
	projectPublisher  *nats.Publisher[*v1.Project]
	buildSubscriber   *nats.Subscriber[*v1.Build]
	clusterSubscriber *nats.Subscriber[*v1.Cluster]
	projectSubscriber *nats.Subscriber[*v1.Project]
}

func NewQueue(js natsgo.JetStreamContext) port.Queue {
	return &Queue{
		buildPublisher:    nats.NewPublisher[*v1.Build](js),
		clusterPublisher:  nats.NewPublisher[*v1.Cluster](js),
		projectPublisher:  nats.NewPublisher[*v1.Project](js),
		buildSubscriber:   nats.NewSubscriber[*v1.Build](js),
		clusterSubscriber: nats.NewSubscriber[*v1.Cluster](js),
		projectSubscriber: nats.NewSubscriber[*v1.Project](js),
	}
}
