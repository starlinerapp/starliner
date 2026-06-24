package http

import (
	"go.uber.org/fx"
	"starliner.app/internal/api/presentation/http/handler"
	"starliner.app/internal/api/presentation/http/middleware"
	"starliner.app/internal/api/presentation/http/sse"
)

var Module = fx.Module(
	"api",
	middleware.Module,
	handler.Module,
	sse.Module,
	fx.Provide(NewServer),
	fx.Invoke(RegisterServer),
)
