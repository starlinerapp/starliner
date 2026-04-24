package handler

import "go.uber.org/fx"

var Module = fx.Module(
	"build-grpc-handlers",
	fx.Provide(NewBuildLogHandler),
)
