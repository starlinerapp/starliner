package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/cluster/application"
	"starliner.app/internal/cluster/conf"
	"starliner.app/internal/cluster/infrastructure/helm"
	sub "starliner.app/internal/cluster/presentation/pubsub"
	cluster "starliner.app/internal/cluster/presentation/queue"
	"starliner.app/internal/core/infrastructure/crypto"
	"starliner.app/internal/core/infrastructure/nats/impl/pubsub"
	"starliner.app/internal/core/infrastructure/nats/impl/queue"
)

func main() {
	fx.New(
		conf.Module,
		crypto.Module,
		helm.Module,
		queue.Module,
		pubsub.Module,
		application.Module,
		cluster.Module,
		sub.Module,
	).Run()
}
