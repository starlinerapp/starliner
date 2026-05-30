package dns

import "go.uber.org/fx"

var Module = fx.Module(
	"dns",
	fx.Provide(
		NewResolver,
	),
)
