package ssh

import "go.uber.org/fx"

var Module = fx.Module(
	"ssh",
	fx.Provide(
		NewSSH,
	),
)
