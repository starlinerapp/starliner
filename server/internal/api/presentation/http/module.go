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
	fx.Provide(NewServer),
	fx.Provide(sse.NewEnvironmentNotificationHub),
	fx.Provide(sse.NewUserNotificationHub),
	fx.Invoke(RegisterServer),
)
