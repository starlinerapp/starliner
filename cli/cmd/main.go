package main

import (
	"go.uber.org/fx"
	"starliner.app/cli/internal/infrastructure/auth"
	"starliner.app/cli/internal/presentation/cli"
)

func main() {
	fx.New(
		fx.NopLogger,
		auth.Module,
		cli.Module,
	).Run()
}
