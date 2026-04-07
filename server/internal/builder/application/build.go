package application

import (
	"context"
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

	tmpDir, commitHash, err := ba.git.CloneRepository(build.GitUrl, build.AccessToken)
	if err != nil {
		log.Printf("failed to clone repository: %v", err)
	}

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Printf("failed to remove directory: %v\n", err)
		}
	}(tmpDir)

	projectDir := filepath.Join(tmpDir, build.RootDirectory)

	imagePath := path.Join(
		ba.cfg.ImageRegistryUrl,
		build.ImageName,
	)
	tag := imagePath + ":" + commitHash

	logs, err := ba.docker.BuildAndPublish(ctx, projectDir, build.DockerfilePath, tag)

	buildStatus := value.BuildStatusSuccess
	if err != nil {
		buildStatus = value.BuildStatusFailed
	}

	err = ba.queue.PublishBuildCompleted(&value.BuildCompleted{
		BuildId:          build.BuildId,
		DeploymentId:     build.DeploymentId,
		ImageRegistryUrl: ba.cfg.ImageRegistryUrl,
		ImageName:        build.ImageName,
		CommitHash:       commitHash,
		Tag:              commitHash,
		Logs:             logs,
		BuildStatus:      buildStatus,
	})
	if err != nil {
		log.Printf("failed to publish: %v", err)
	}
}
