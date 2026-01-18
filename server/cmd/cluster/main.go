package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/cluster/application"
	"starliner.app/internal/cluster/infrastructure/helm"
	cluster "starliner.app/internal/cluster/presentation/queue"
	"starliner.app/internal/core/conf"
	"starliner.app/internal/core/domain/repository"
	"starliner.app/internal/core/infrastructure/crypto"
	"starliner.app/internal/core/infrastructure/nats/impl/queue"
	"starliner.app/internal/core/infrastructure/postgres"
)

func main() {
	fx.New(
		conf.Module,
		crypto.Module,
		postgres.Module,
		helm.Module,
		queue.Module,
		repository.Module,
		application.Module,
		cluster.Module,
	).Run()
}
