package docker

import (
	"context"
	"errors"
	"fmt"
	"github.com/moby/buildkit/client"
	gateway "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/util/progress/progresswriter"
	"github.com/tonistiigi/fsutil"
	"os"
	"path/filepath"
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
	bkClient, err := client.New(ctx, "tcp://buildkit:1234")
	if err != nil {
		return "", fmt.Errorf("failed to connect to buildkit at %s: %w", "tcp://buildkit:1234", err)
	}
	defer func(bkClient *client.Client) {
		_ = bkClient.Close()
	}(bkClient)

	absProjectDir, err := filepath.Abs(projectDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve project dir: %w", err)
	}
	absDockerfilePath, err := filepath.Abs(dockerfilePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve dockerfile path: %w", err)
	}

	pw, err := progresswriter.NewPrinter(ctx, os.Stdout, "auto")
	if err != nil {
		return "", fmt.Errorf("failed to create progress writer: %w", err)
	}

	attachable := []session.Attachable{
		authprovider.NewDockerAuthProvider(authprovider.DockerAuthProviderConfig{}),
	}

	contextFS, err := fsutil.NewFS(absProjectDir)
	if err != nil {
		return "", fmt.Errorf("failed to create context FS: %w", err)
	}

	dockerfileFS, err := fsutil.NewFS(filepath.Dir(absDockerfilePath))
	if err != nil {
		return "", fmt.Errorf("failed to create dockerfile FS: %w", err)
	}

	res, err := bkClient.Build(ctx, client.SolveOpt{
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
			"dockerfile": dockerfileFS,
		},
		FrontendAttrs: map[string]string{
			"filename": filepath.Base(absDockerfilePath),
		},
		Session: attachable,
	}, "", func(ctx context.Context, c gateway.Client) (*gateway.Result, error) {
		return c.Solve(ctx, gateway.SolveRequest{
			Frontend: "dockerfile.v0",
		})
	}, nil)

	<-pw.Done()
	if err := pw.Err(); err != nil && !errors.Is(err, context.Canceled) {
		return "", fmt.Errorf("build progress error: %w", err)
	}

	if err != nil {
		return "", fmt.Errorf("build and push failed: %w", err)
	}

	digest, ok := res.ExporterResponse["containerimage.digest"]
	if !ok || digest == "" {
		return "", fmt.Errorf("no digest in exporter response")
	}

	return digest, nil
}
