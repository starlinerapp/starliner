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
	cfg    *conf.Config
	git    port.Git
	docker port.Docker
	queue  port.Queue
}

func NewBuildApplication(
	cfg *conf.Config,
	git port.Git,
	docker port.Docker,
	queue port.Queue,
) *BuildApplication {
	return &BuildApplication{
		cfg:    cfg,
		git:    git,
		docker: docker,
		queue:  queue,
	}
}

func (ba *BuildApplication) HandleBuildTriggered(build *value.TriggerBuild) {
	ctx := context.Background()

	publishCompleted := func(commitHash, tag *string, logs string, status value.BuildStatus) {
		if err := ba.queue.PublishBuildCompleted(&value.BuildCompleted{
			BuildId:          build.BuildId,
			DeploymentId:     build.DeploymentId,
			ImageRegistryUrl: ba.cfg.ImageRegistryUrl,
			ImageName:        build.ImageName,
			CommitHash:       commitHash,
			Tag:              commitHash,
			Logs:             logs,
			BuildStatus:      status,
		}); err != nil {
			log.Printf("failed to publish: %v", err)
		}
	}

	tmpDir, commitHash, err := ba.git.CloneRepository(build.GitUrl, build.AccessToken)
	if err != nil {
		publishCompleted(
			nil,
			nil,
			fmt.Sprintf("failed to clone repository: %v", err),
			value.BuildStatusFailed,
		)
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

	logs, err := ba.docker.BuildAndPublish(ctx, projectDir, build.DockerfilePath, tag, build.Args)

	status := value.BuildStatusSuccess
	if err != nil {
		status = value.BuildStatusFailed
	}

	publishCompleted(&commitHash, &tag, logs, status)
}
