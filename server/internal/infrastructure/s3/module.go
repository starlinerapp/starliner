package s3

import (
	"context"
	"go.uber.org/fx"
	"starliner.app/internal/domain/port"
)

var Module = fx.Module(
	"objectstore",
	fx.Provide(
		Connect,
		NewS3Client,
	),
	fx.Invoke(func(c port.ObjectStore, lc fx.Lifecycle) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return c.CreateBuckets(ctx)
			},
		})
	}),
)
