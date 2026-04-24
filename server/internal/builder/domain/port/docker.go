package port

import (
	"context"

	"starliner.app/internal/core/domain/value"
)

type Docker interface {
	BuildAndPublish(
		ctx context.Context,
		buildId int64,
		projectDir,
		dockerfilePath string,
		imageTag string,
		args []*value.Arg,
	) (string, error)
}
