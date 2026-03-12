package docker

import (
	"bytes"
	"context"
	"dagger.io/dagger"
	"regexp"
	"starliner.app/internal/builder/domain/port"
)

type Docker struct {
}

func NewDocker() port.Docker {
	return &Docker{}
}

func (c *Docker) BuildAndPublish(
	ctx context.Context,
	projectDir string,
	dockerfilePath string,
	imageTag string,
) (string, error) {
	var logBuffer bytes.Buffer

	client, err := dagger.Connect(
		ctx,
		dagger.WithLogOutput(&logBuffer),
	)
	if err != nil {
		return "", err
	}
	defer func(client *dagger.Client) {
		_ = client.Close()
	}(client)

	buildContainer := client.Host().
		Directory(projectDir).
		DockerBuild(dagger.DirectoryDockerBuildOpts{
			Dockerfile: dockerfilePath,
			Platform:   "linux/amd64",
		})

	_, err = buildContainer.Publish(ctx, imageTag)
	if err != nil {
		return logBuffer.String(), err
	}

	return stripANSI(logBuffer.String()), nil
}

func stripANSI(s string) string {
	var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}
