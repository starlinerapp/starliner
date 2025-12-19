package builder

import (
	"context"
	"go.uber.org/fx"
	"log"
	"starliner.app/pkg/objectstore"
	v1 "starliner.app/pkg/proto/v1"
	"starliner.app/pkg/queue"
)

type Orchestrator struct {
	objectstore              *objectstore.S3Client
	projectCreatedSubscriber *queue.Subscriber[*v1.Project]
}

func NewOrchestrator(
	objectstore *objectstore.S3Client,
	projectCreatedSubscriber *queue.Subscriber[*v1.Project],
) *Orchestrator {
	return &Orchestrator{
		objectstore:              objectstore,
		projectCreatedSubscriber: projectCreatedSubscriber,
	}
}

func (o *Orchestrator) Start() error {
	go func() {
		err := o.projectCreatedSubscriber.Subscribe(queue.ProjectCreated, "orchestrator", o.handleProjectCreated)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()

	return nil
}

func (o *Orchestrator) handleProjectCreated(project *v1.Project) {
	log.Printf("project name: %v\n", project.Name)
}

func RegisterOrchestrator(lc fx.Lifecycle, o *Orchestrator) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return o.Start()
		},
	})
}
