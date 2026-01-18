package dagger

import "go.uber.org/fx"

var Module = fx.Module(
	"dagger",
	fx.Provide(
		NewDaggerClient,
	),
)
