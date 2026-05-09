package main

import (
	"go.uber.org/fx"
	"starliner.app/cli/internal/presentation/cli"
)

func main() {
	fx.New(
		fx.NopLogger,
		cli.Module,
	).Run()
}
