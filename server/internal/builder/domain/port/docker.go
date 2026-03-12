package port

import (
	"context"
)

type Docker interface {
	BuildAndPublish(
		ctx context.Context,
		projectDir,
		dockerfilePath string,
		imageTag string,
	) (string, error)
}
