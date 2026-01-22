package application

import (
	"github.com/google/uuid"
	"log"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/core/domain/value"
)

type BuildApplication struct {
	queue port.Queue
}

func NewBuildApplication(
	queue port.Queue,
) *BuildApplication {
	return &BuildApplication{
		queue: queue,
	}
}

func (ba *BuildApplication) TriggerBuild() error {
	buildId := uuid.New().String()
	err := ba.queue.PublishBuildTriggered(&value.TriggerBuild{
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
