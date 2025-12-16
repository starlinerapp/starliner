package main

import (
	"go.uber.org/fx"
	"starliner.app/pkg/api"
	"starliner.app/pkg/config"
)

// @title Starliner API
// @version 1.0
// @Param X-User-ID header string true "User ID"
// @securityDefinitions.basic BasicAuth
func main() {
	fx.New(
		config.Module,
		api.Module,
	).Run()
}
