package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/application"
	"starliner.app/internal/conf"
	"starliner.app/internal/domain/repository"
	"starliner.app/internal/domain/service"
	"starliner.app/internal/infrastructure/crypto"
	"starliner.app/internal/infrastructure/dagger"
	"starliner.app/internal/infrastructure/db"
	"starliner.app/internal/infrastructure/objectstore"
	"starliner.app/internal/infrastructure/queue"
	"starliner.app/internal/presentation/queue/builder"
	"starliner.app/internal/presentation/queue/cluster"
	"starliner.app/internal/presentation/queue/provisioner"
)

func main() {
	fx.New(
		conf.Module,
		db.Module,
		repository.Module,
		queue.Module,
		objectstore.Module,
		dagger.Module,
		crypto.Module,
		application.Module,
		service.Module,
		builder.Module,
		provisioner.Module,
		cluster.Module,
	).Run()
}
