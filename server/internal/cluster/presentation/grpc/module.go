package grpc

import (
	"go.uber.org/fx"
	"starliner.app/internal/cluster/presentation/grpc/handler"
)

var Module = fx.Module(
	"grpc",
	handler.Module,
	fx.Provide(NewServer),
	fx.Invoke(RegisterServer),
)
