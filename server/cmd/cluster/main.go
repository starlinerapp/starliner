package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/cluster/application"
	"starliner.app/internal/cluster/conf"
	"starliner.app/internal/cluster/infrastructure/helm"
	"starliner.app/internal/cluster/infrastructure/nats/impl/pubsub"
	"starliner.app/internal/cluster/infrastructure/nats/impl/queue"
	sub "starliner.app/internal/cluster/presentation/pubsub"
	clusterqueue "starliner.app/internal/cluster/presentation/queue"
	"starliner.app/internal/core/infrastructure/crypto"
)

func main() {
	fx.New(
		conf.Module,
		crypto.Module,
		helm.Module,
		queue.Module,
		pubsub.Module,
		application.Module,
		clusterqueue.Module,
		sub.Module,
	).Run()
}
