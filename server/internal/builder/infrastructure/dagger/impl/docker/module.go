package docker

import "go.uber.org/fx"

var Module = fx.Module(
	"docker",
	fx.Provide(
		NewDocker,
	),
)
