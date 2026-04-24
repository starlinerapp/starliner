package application

import (
	"go.uber.org/fx"
	"starliner.app/internal/builder/domain/port"
)

var Module = fx.Module(
	"application",
	fx.Provide(
		NewBuildLogApplication,
		func(ba *BuildLogApplication) port.LogPublisher {
			return ba
		},
		NewBuildApplication,
	),
)
