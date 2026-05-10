package wg

import "go.uber.org/fx"

var Module = fx.Module(
	"wg",
	fx.Provide(
		NewClient,
	),
)
