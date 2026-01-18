package pulumi

import "go.uber.org/fx"

var Module = fx.Module(
	"pulumi",
	fx.Provide(
		NewProvision,
	),
)
