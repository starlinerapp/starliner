package main

import (
	"go.uber.org/fx"
	"starliner.app/pkg/config"
)

func main() {
	fx.New(
		config.Module,
	).Run()
}
