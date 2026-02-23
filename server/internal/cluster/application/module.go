package application

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"application",
	fx.Provide(
		NewImageApplication,
		NewDatabaseApplication,
		NewIngressApplication,
		NewStatusApplication,
	),
)
