package redis

import "go.uber.org/fx"

var Module = fx.Module(
	"redis",
	fx.Provide(
		Connect,
		NewClient,
	),
)
