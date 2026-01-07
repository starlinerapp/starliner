package docker

import (
	"context"
	"dagger.io/dagger"
	"starliner.app/internal/domain/port"
)

type Docker struct {
	daggerClient *dagger.Client
}

func NewDocker(
	daggerClient *dagger.Client,
) port.Docker {
	return &Docker{
		daggerClient: daggerClient,
	}
}

func (c *Docker) BuildAndPublish(ctx context.Context, projectDir string, dockerfilePath string, imageTag string) error {
	buildContainer := c.daggerClient.Host().
		Directory(projectDir).
		DockerBuild(dagger.DirectoryDockerBuildOpts{Dockerfile: dockerfilePath})

	_, err := buildContainer.Publish(ctx, imageTag)
	return err
}
