package application

import (
	"context"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"starliner.app/internal/conf"
	"starliner.app/internal/domain/port"
	"starliner.app/internal/domain/value"
	"strings"
)

type BuildApplication struct {
	cfg         *conf.Config
	docker      port.Docker
	objectstore port.ObjectStore
	queue       port.Queue
}

func NewBuildApplication(
	cfg *conf.Config,
	docker port.Docker,
	objectstore port.ObjectStore,
	queue port.Queue,
) *BuildApplication {
	return &BuildApplication{
		cfg:         cfg,
		objectstore: objectstore,
		docker:      docker,
		queue:       queue,
	}
}

func (ba *BuildApplication) TriggerBuild() error {
	buildId := uuid.New().String()
	err := ba.queue.PublishBuildTriggered(&value.BuildMessage{
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

func (ba *BuildApplication) HandleBuildTriggered(build *value.BuildMessage) {
	ctx := context.Background()

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

	body, err := ba.objectstore.GetObject(ctx, build.S3Key)
	if err != nil {
		log.Printf("failed to get file from S3: %v", err)
		return
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}(body)

	tarFileName := filepath.Base(build.S3Key)
	tarPath := filepath.Join(workDir, tarFileName)

	f, err := os.Create(tarPath)
	if err != nil {
		log.Printf("failed to create tarball: %v", err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("failed to close tarball: %v", err)
		}
	}(f)

	if _, err := io.Copy(f, body); err != nil {
		log.Printf("failed to write tarball: %v", err)
		return
	}

	cmd := exec.Command("tar", "-xzf", tarPath, "-C", workDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("failed to extract tarball: %v, output: %s", err, string(out))
		return
	}

	extractDirName := strings.TrimSuffix(tarFileName, ".tgz")
	extractDir := filepath.Join(workDir, extractDirName)

	imagePath := path.Join(
		ba.cfg.ImageRegistryUrl,
		build.Organization,
		build.Project,
		build.Service,
	)
	tag := imagePath + ":latest"
	projectDir := filepath.Join(extractDir, build.RootDirectory)

	err = ba.docker.BuildAndPublish(ctx, projectDir, build.DockerfilePath, tag)
	if err != nil {
		log.Printf("failed to build docker image: %v", err)
		return
	}
}
