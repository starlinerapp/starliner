package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/config"
	"starliner.app/internal/core/builder"
	"starliner.app/internal/core/cluster"
	"starliner.app/internal/core/provisioner"
	"starliner.app/internal/infrastructure/db"
	"starliner.app/internal/infrastructure/queue"
	"starliner.app/internal/repository"
	"starliner.app/internal/service"
)

func main() {
	fx.New(
		config.Module,
		db.Module,
		repository.Module,
		service.Module,
		queue.Module,
		builder.Module,
		provisioner.Module,
		cluster.Module,
	).Run()
}
