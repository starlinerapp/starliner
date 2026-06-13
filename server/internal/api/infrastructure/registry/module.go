package registry

import "go.uber.org/fx"

var Module = fx.Module(
	"registry",
	fx.Provide(
		NewClient,
	),
)
