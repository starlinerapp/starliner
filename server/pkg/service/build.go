package service

import (
	"github.com/google/uuid"
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
	buildId := uuid.New().String()
	err := bs.buildPublisher.Publish(queue.BuildTriggered, buildId, &v1.Build{
		Id:             buildId,
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
