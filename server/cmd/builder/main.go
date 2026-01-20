package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/builder/application"
	"starliner.app/internal/builder/conf"
	"starliner.app/internal/builder/infrastructure/dagger"
	"starliner.app/internal/builder/infrastructure/dagger/impl/docker"
	builder "starliner.app/internal/builder/presentation/queue"
	"starliner.app/internal/core/infrastructure/nats/impl/queue"
	"starliner.app/internal/core/infrastructure/s3"
)

func main() {
	fx.New(
		conf.Module,
		s3.Module,
		queue.Module,
		dagger.Module,
		docker.Module,
		application.Module,
		builder.Module,
	).Run()
}
