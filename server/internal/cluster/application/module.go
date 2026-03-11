package application

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"application",
	fx.Provide(
		NewDeploymentApplication,
		NewImageApplication,
		NewDatabaseApplication,
		NewIngressApplication,
		NewStatusApplication,
		NewLogsApplication,
		NewTTYApplication,
	),
)
