package handler

import "go.uber.org/fx"

var Module = fx.Module(
	"api-grpc-handlers",
	fx.Provide(
		NewRunnerHandler,
	),
)
