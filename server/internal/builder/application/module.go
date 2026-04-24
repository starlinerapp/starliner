package application

import (
	"go.uber.org/fx"
	"starliner.app/internal/builder/domain/port"
)

var Module = fx.Module(
	"application",
	fx.Provide(
		NewBuildLogApplication,
		func(a *BuildLogApplication) port.LogPublisher {
			return a
		},
		NewBuildApplication,
	),
)
