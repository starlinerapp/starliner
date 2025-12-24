package service

import (
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

func (bs *BuildService) TriggerBuild() error {
	err := bs.buildPublisher.Publish(queue.BuildTriggered, &v1.Build{
		Organization:   "starliner",
		Project:        "example",
		Service:        "client",
		S3Key:          "monorepo-example.tgz",
		RootDirectory:  "./client",
		DockerfilePath: "Dockerfile",
	})

	if err != nil {
		log.Printf("error publishing: %v", err)
	}

	return nil
}
