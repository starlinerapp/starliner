package middleware

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"middleware",
	fx.Provide(
		NewUserMiddleware,
		NewBasicAuthMiddleware,
	),
)
