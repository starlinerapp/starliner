package ansible

import "go.uber.org/fx"

var Module = fx.Module(
	"install",
	fx.Provide(NewInstall),
)
