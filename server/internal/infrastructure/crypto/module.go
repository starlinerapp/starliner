package crypto

import "go.uber.org/fx"

var Module = fx.Module(
	"crypto",
	fx.Provide(
		NewCrypto,
	),
)
