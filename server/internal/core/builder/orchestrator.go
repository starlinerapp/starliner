package builder

import (
	"context"
	"dagger.io/dagger"
	"go.uber.org/fx"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"starliner.app/internal/config"
	"starliner.app/internal/infrastructure/objectstore"
	"starliner.app/internal/infrastructure/queue"
	"starliner.app/internal/infrastructure/queue/proto/v1"
	"strings"
)

type Orchestrator struct {
	cfg             *config.Config
	objectstore     *objectstore.S3Client
	buildSubscriber *queue.Subscriber[*v1.Build]
	daggerClient    *dagger.Client
}

func RegisterOrchestrator(lc fx.Lifecycle, o *Orchestrator) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return o.Start()
		},
	})
}

func NewOrchestrator(
	cfg *config.Config,
	objectstore *objectstore.S3Client,
	buildSubscriber *queue.Subscriber[*v1.Build],
	daggerClient *dagger.Client,
) *Orchestrator {
	return &Orchestrator{
		cfg:             cfg,
		objectstore:     objectstore,
		buildSubscriber: buildSubscriber,
		daggerClient:    daggerClient,
	}
}

func (o *Orchestrator) Start() error {
	go func() {
		err := o.buildSubscriber.Subscribe(queue.BuildTriggered, "*", "buildTriggered", o.handleBuildTriggered)
		if err != nil {
			log.Fatalf("failed to subscribe to queue: %v", err)
		}
	}()

	return nil
}

func (o *Orchestrator) handleBuildTriggered(build *v1.Build) {
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

	body, err := o.objectstore.GetObject(ctx, build.S3Key)
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
		o.cfg.ImageRegistryUrl,
		build.Organization,
		build.Project,
		build.Service,
	)

	projectDir := filepath.Join(extractDir, build.RootDirectory)
	buildContainer := o.daggerClient.Host().
		Directory(projectDir).
		DockerBuild(dagger.DirectoryDockerBuildOpts{Dockerfile: build.DockerfilePath})

	_, err = buildContainer.Publish(ctx, imagePath+":latest")
	if err != nil {
		log.Printf("failed to build docker image: %v", err)
		return
	}
}
