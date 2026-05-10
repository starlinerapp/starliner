package k3d

import "go.uber.org/fx"

var Module = fx.Module(
	"k3d",
	fx.Provide(
		NewClient,
	),
)
