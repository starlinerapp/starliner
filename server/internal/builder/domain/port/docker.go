package port

import (
	"context"

	"starliner.app/internal/core/domain/value"
)

type Docker interface {
	BuildAndPublish(
		ctx context.Context,
		projectDir,
		dockerfilePath string,
		imageTag string,
		args []*value.Arg,
	) (string, error)
}
