package helm

import "go.uber.org/fx"

var Module = fx.Module(
	"helm",
	fx.Provide(
		NewDeploy,
	),
)
