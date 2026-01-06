package application

import (
	"context"
	"dagger.io/dagger"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"starliner.app/internal/conf"
	"starliner.app/internal/infrastructure/objectstore"
	"starliner.app/internal/infrastructure/queue"
	"starliner.app/internal/infrastructure/queue/proto/v1"
	"strings"
)

type BuildApplication struct {
	cfg            *conf.Config
	objectstore    *objectstore.S3Client
	daggerClient   *dagger.Client
	buildPublisher *queue.Publisher[*v1.Build]
}

func NewBuildApplication(
	cfg *conf.Config,
	objectstore *objectstore.S3Client,
	daggerClient *dagger.Client,
	buildPublisher *queue.Publisher[*v1.Build],
) *BuildApplication {
	return &BuildApplication{
		cfg:            cfg,
		objectstore:    objectstore,
		daggerClient:   daggerClient,
		buildPublisher: buildPublisher,
	}
}

func (ba *BuildApplication) TriggerBuild() error {
	buildId := uuid.New().String()
	err := ba.buildPublisher.Publish(queue.BuildTriggered, buildId, &v1.Build{
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

func (ba *BuildApplication) HandleBuildTriggered(build *v1.Build) {
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

	projectDir := filepath.Join(extractDir, build.RootDirectory)
	buildContainer := ba.daggerClient.Host().
		Directory(projectDir).
		DockerBuild(dagger.DirectoryDockerBuildOpts{Dockerfile: build.DockerfilePath})

	_, err = buildContainer.Publish(ctx, imagePath+":latest")
	if err != nil {
		log.Printf("failed to build docker image: %v", err)
		return
	}
}
