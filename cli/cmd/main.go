package main

import (
	"go.uber.org/fx"
	"starliner.app/cli/internal/application"
	"starliner.app/cli/internal/conf"
	"starliner.app/cli/internal/infrastructure/auth"
	"starliner.app/cli/internal/infrastructure/k3d"
	"starliner.app/cli/internal/infrastructure/wg"
	"starliner.app/cli/internal/presentation/cli"
)

func main() {
	fx.New(
		fx.NopLogger,
		conf.Module,
		auth.Module,
		wg.Module,
		k3d.Module,
		application.Module,
		cli.Module,
	).Run()
}
