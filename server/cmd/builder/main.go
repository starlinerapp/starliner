package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/builder/application"
	"starliner.app/internal/builder/conf"
	docker "starliner.app/internal/builder/infrastructure/dagger"
	"starliner.app/internal/builder/infrastructure/git"
	"starliner.app/internal/builder/infrastructure/nats/impl/queue"
	builderqueue "starliner.app/internal/builder/presentation/queue"
	"starliner.app/internal/core/infrastructure/s3"
)

func main() {
	fx.New(
		conf.Module,
		s3.Module,
		queue.Module,
		git.Module,
		docker.Module,
		application.Module,
		builderqueue.Module,
	).Run()
}
