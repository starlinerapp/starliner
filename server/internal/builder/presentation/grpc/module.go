package grpc

import (
	"go.uber.org/fx"
	"starliner.app/internal/builder/presentation/grpc/handler"
)

var Module = fx.Module(
	"build-grpc",
	handler.Module,
	fx.Provide(
		NewServer,
	),
	fx.Invoke(RegisterServer),
)
