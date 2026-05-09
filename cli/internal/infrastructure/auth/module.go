package auth

import (
	"go.uber.org/fx"
	"starliner.app/cli/internal/domain/port"
)

var Module = fx.Module(
	"auth",
	fx.Provide(
		fx.Annotate(
			NewClient,
			fx.As(new(port.AuthClient)),
		),
	),
)
