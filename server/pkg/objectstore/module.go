package objectstore

import (
	"context"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"objectstore",
	fx.Provide(
		Connect,
		NewS3Client,
	),
	fx.Invoke(func(c *S3Client, lc fx.Lifecycle) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return c.CreateBucket(ctx)
			},
		})
	}),
)
