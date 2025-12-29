package main

import (
	"go.uber.org/fx"
	"starliner.app/pkg/builder"
	"starliner.app/pkg/cluster"
	"starliner.app/pkg/config"
	"starliner.app/pkg/queue"
)

func main() {
	fx.New(
		config.Module,
		queue.Module,
		builder.Module,
		cluster.Module,
	).Run()
}
