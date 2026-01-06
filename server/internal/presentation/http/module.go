package http

import (
	"go.uber.org/fx"
	"starliner.app/internal/presentation/http/handler"
	"starliner.app/internal/presentation/http/middleware"
)

var Module = fx.Module(
	"api",
	middleware.Module,
	handler.Module,
	fx.Provide(NewServer),
	fx.Invoke(RegisterServer),
)
