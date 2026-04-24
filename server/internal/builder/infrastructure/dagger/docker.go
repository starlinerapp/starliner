package docker

import (
	"bytes"
	"context"
	"regexp"

	"dagger.io/dagger"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
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
	args []*value.Arg,
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

	var buildArgs []dagger.BuildArg
	for _, a := range args {
		if a != nil {
			buildArgs = append(buildArgs, dagger.BuildArg{
				Name:  a.Name,
				Value: a.Value,
			})
		}
	}

	buildContainer := client.Host().
		Directory(projectDir).
		DockerBuild(dagger.DirectoryDockerBuildOpts{
			Dockerfile: dockerfilePath,
			Platform:   "linux/amd64",
			BuildArgs:  buildArgs,
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
