package builder

import (
	"context"
	"go.uber.org/fx"
	"log"
	"os"
	"starliner.app/pkg/objectstore"
	v1 "starliner.app/pkg/proto/v1"
	"starliner.app/pkg/queue"
	"starliner.app/pkg/utils"
)

type Orchestrator struct {
	objectstore     *objectstore.S3Client
	buildSubscriber *queue.Subscriber[*v1.Build]
}

func NewOrchestrator(
	objectstore *objectstore.S3Client,
	buildSubscriber *queue.Subscriber[*v1.Build],
) *Orchestrator {
	return &Orchestrator{
		objectstore:     objectstore,
		buildSubscriber: buildSubscriber,
	}
}

func (o *Orchestrator) Start(ctx context.Context) error {
	go func() {
		err := o.buildSubscriber.Subscribe(ctx, queue.BuildTriggered, "buildTriggered", o.handleBuildTriggered)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()

	return nil
}

func (o *Orchestrator) handleBuildTriggered(ctx context.Context, build *v1.Build) {
	workDir, err := os.MkdirTemp("", "build-*")
	if err != nil {
		log.Printf("failed to create temp dir: %v", err)
		return
	}
	defer func() {
		if err := os.RemoveAll(workDir); err != nil {
			log.Printf("failed to cleanup %s: %v", workDir, err)
		}
	}()

	zip, err := o.objectstore.GetFile(ctx, build.S3Key)
	if err != nil {
		log.Printf("failed to get file from S3: %v", err)
		return
	}

	err = utils.Unzip(zip, workDir)
	if err != nil {
		log.Printf("failed to unzip file: %v", err)
		return
	}
}

func RegisterOrchestrator(lc fx.Lifecycle, o *Orchestrator) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return o.Start(ctx)
		},
	})
}
