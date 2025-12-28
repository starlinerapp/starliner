package main

import (
	"go.uber.org/fx"
	"starliner.app/pkg/api"
	"starliner.app/pkg/config"
)

// @title Starliner API
// @version 1.0
// @securityDefinitions.basic BasicAuth
func main() {
	fx.New(
		config.Module,
		api.Module,
	).Run()
}
