package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/config"
	"starliner.app/internal/core/builder"
	"starliner.app/internal/core/cluster"
	"starliner.app/internal/infrastructure/queue"
)

func main() {
	fx.New(
		config.Module,
		queue.Module,
		builder.Module,
		cluster.Module,
	).Run()
}
