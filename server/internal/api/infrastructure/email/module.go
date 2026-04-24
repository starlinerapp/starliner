package email

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"email",
	fx.Provide(
		NewClient,
	),
)
