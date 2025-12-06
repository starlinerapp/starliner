package main

import (
	"log"
	"starliner.app/pkg/config"
	"starliner.app/pkg/di"
)

// @title   Starliner API
// @version 1.0
func main() {
	cfg, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	server, diErr := di.InitializeAPI(cfg)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	}

	server.Start()
}
