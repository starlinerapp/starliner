package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"starliner.app/internal/builder/conf"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
)

type BuildApplication struct {
	cfg          *conf.Config
	git          port.Git
	docker       port.Docker
	queue        port.Queue
	logPublisher port.LogPublisher
}

func NewBuildApplication(
	cfg *conf.Config,
	git port.Git,
	docker port.Docker,
	queue port.Queue,
	logPublisher port.LogPublisher,
) *BuildApplication {
	return &BuildApplication{
		cfg:          cfg,
		git:          git,
		docker:       docker,
		queue:        queue,
		logPublisher: logPublisher,
	}
}

func (ba *BuildApplication) HandleBuildTriggered(build *value.TriggerBuild) {
	ctx := context.Background()

	publishLogLine := func(line string) {
		if ba.logPublisher == nil {
			return
		}
		if err := ba.logPublisher.PublishLogChunk(build.BuildId, []byte(line)); err != nil {
			log.Printf("failed to publish log chunk: %v", err)
		}
	}

	// Always emit an end marker before BuildCompleted so that any active
	// log subscribers can release their connection.
	defer func() {
		if ba.logPublisher == nil {
			return
		}
		if err := ba.logPublisher.PublishLogEnd(build.BuildId); err != nil {
			log.Printf("failed to publish log end: %v", err)
		}
	}()

	publishCompleted := func(commitHash, tag *string, imageName *string, logs string, status value.BuildStatus) {
		if err := ba.queue.PublishBuildCompleted(&value.BuildCompleted{
			BuildId:          build.BuildId,
			DeploymentId:     build.DeploymentId,
			ImageRegistryUrl: ba.cfg.ImageRegistryUrl,
			ImageName:        imageName,
			CommitHash:       commitHash,
			Tag:              commitHash,
			Logs:             logs,
			BuildStatus:      status,
		}); err != nil {
			log.Printf("failed to publish: %v", err)
		}
	}

	tmpDir, commitHash, err := ba.git.CloneRepository(build.GitUrl, build.BranchName, build.AccessToken)
	if err != nil {
		msg := fmt.Sprintf("failed to clone repository: %v", err)
		publishLogLine(msg + "\n")
		publishCompleted(nil, nil, nil, msg, value.BuildStatusFailed)
		return
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Printf("failed to remove directory: %v", err)
		}
	}()

	projectDir := filepath.Join(tmpDir, build.RootDirectory)
	imagePath := path.Join(ba.cfg.ImageRegistryUrl, build.ImageName)
	tag := imagePath + ":" + commitHash

	logs, err := ba.docker.BuildAndPublish(ctx, build.BuildId, projectDir, build.DockerfilePath, tag, build.Args)

	status := value.BuildStatusSuccess
	if err != nil {
		msg := fmt.Sprintf("✗ %v", err)
		publishLogLine(msg + "\n")
		logs += msg + "\n"
		status = value.BuildStatusFailed
	}

	publishCompleted(&commitHash, &tag, &imagePath, logs, status)
}
