package application

import (
	"github.com/google/uuid"
	"log"
	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/core/domain/value"
)

type BuildApplication struct {
	queue corePort.Queue
}

func NewBuildApplication(
	queue corePort.Queue,
) *BuildApplication {
	return &BuildApplication{
		queue: queue,
	}
}

func (ba *BuildApplication) TriggerBuild() error {
	buildId := uuid.New().String()
	err := ba.queue.PublishBuildTriggered(&value.Build{
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
