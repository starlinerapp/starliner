package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/conf"
	"starliner.app/internal/presentation/http"
)

// @title Starliner API
// @version 1.0
// @securityDefinitions.basic BasicAuth
func main() {
	fx.New(
		conf.Module,
		http.Module,
	).Run()
}
