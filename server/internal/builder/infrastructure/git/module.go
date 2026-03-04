package git

import "go.uber.org/fx"

var Module = fx.Module(
	"git",
	fx.Provide(
		NewGit,
	),
)
