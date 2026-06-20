package grpc

import (
	"go.uber.org/fx"
	"starliner.app/internal/api/presentation/grpc/handler"
)

var Module = fx.Module(
	"api-grpc",
	handler.Module,
	fx.Provide(
		NewServer,
	),
	fx.Invoke(RegisterServer),
)
