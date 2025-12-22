package service

import (
	"context"
	"log"
	v1 "starliner.app/pkg/proto/v1"
	"starliner.app/pkg/queue"
)

type BuildService struct {
	buildPublisher *queue.Publisher[*v1.Build]
}

func NewBuildService(buildPublisher *queue.Publisher[*v1.Build]) *BuildService {
	return &BuildService{
		buildPublisher: buildPublisher,
	}
}

func (bs *BuildService) TriggerBuild(ctx context.Context) error {
	err := bs.buildPublisher.Publish(queue.BuildTriggered, &v1.Build{
		S3Key:          "example-project.zip",
		DockerfilePath: ".",
	})

	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}
