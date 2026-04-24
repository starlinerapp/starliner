package docker

import (
	"context"
	"errors"
	"fmt"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/tonistiigi/fsutil"
	"log"
	"os"
	"path/filepath"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
	"strings"
	"sync"
)

type Docker struct{}

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
	bkClient, err := client.New(ctx, "tcp://buildkit:1234")
	if err != nil {
		return "", fmt.Errorf("failed to connect to buildkit: %w", err)
	}
	defer func(bkClient *client.Client) {
		if err := bkClient.Close(); err != nil {
			log.Printf("failed to close buildkit client: %v", err)
		}
	}(bkClient)

	absProjectDir, err := filepath.Abs(projectDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve project dir: %w", err)
	}

	dockerfileRelPath := dockerfilePath
	if filepath.IsAbs(dockerfilePath) {
		dockerfileRelPath, err = filepath.Rel(absProjectDir, dockerfilePath)
		if err != nil {
			return "", fmt.Errorf("failed to compute relative dockerfile path: %w", err)
		}
	}

	dockerfileRelPath = filepath.Clean(dockerfileRelPath)

	if dockerfileRelPath == "." {
		return "", fmt.Errorf("dockerfile path points to the context directory, not a file")
	}

	if filepath.IsAbs(dockerfileRelPath) ||
		dockerfileRelPath == ".." ||
		strings.HasPrefix(dockerfileRelPath, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf(
			"dockerfile must be inside build context: projectDir=%s dockerfilePath=%s resolvedDockerfile=%s",
			absProjectDir,
			dockerfilePath,
			dockerfileRelPath,
		)
	}

	contextFS, err := fsutil.NewFS(absProjectDir)
	if err != nil {
		return "", fmt.Errorf("failed to create context FS: %w", err)
	}

	frontendAttrs := map[string]string{
		"filename": filepath.ToSlash(dockerfileRelPath),
	}

	for _, a := range args {
		if a == nil {
			continue
		}
		frontendAttrs["build-arg:"+a.Name] = a.Value
	}

	attachable := []session.Attachable{
		authprovider.NewDockerAuthProvider(authprovider.DockerAuthProviderConfig{}),
	}

	statusCh := make(chan *client.SolveStatus)
	doneCh := make(chan struct{})

	var (
		logs strings.Builder
		mu   sync.Mutex
	)

	appendLog := func(format string, values ...any) {
		line := fmt.Sprintf(format, values...)

		mu.Lock()
		defer mu.Unlock()

		logs.WriteString(line)
	}

	go func() {
		defer close(doneCh)

		for status := range statusCh {
			for _, log := range status.Logs {
				line := string(log.Data)

				_, err := fmt.Fprint(os.Stdout, line)
				if err != nil {
					return
				}
				appendLog("%s", line)
			}

			for _, s := range status.Statuses {
				if s.Completed != nil {
					line := fmt.Sprintf("✓ %s\n", s.ID)
					_, _ = fmt.Fprint(os.Stdout, line)

					appendLog("%s", line)
				}
			}
		}
	}()

	_, buildErr := bkClient.Solve(
		ctx,
		nil,
		client.SolveOpt{
			Frontend:      "dockerfile.v0",
			FrontendAttrs: frontendAttrs,
			Exports: []client.ExportEntry{
				{
					Type: client.ExporterImage,
					Attrs: map[string]string{
						"name": imageTag,
						"push": "true",
					},
				},
			},
			LocalMounts: map[string]fsutil.FS{
				"context":    contextFS,
				"dockerfile": contextFS,
			},
			Session: attachable,
		},
		statusCh,
	)

	<-doneCh

	mu.Lock()
	logOutput := logs.String()
	mu.Unlock()

	if buildErr != nil {
		if errors.Is(buildErr, context.Canceled) {
			return logOutput, buildErr
		}
		return logOutput, fmt.Errorf("build and push failed: %w", buildErr)
	}
	return logOutput, nil
}
