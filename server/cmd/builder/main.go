package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/builder/application"
	"starliner.app/internal/builder/conf"
	docker "starliner.app/internal/builder/infrastructure/buildkit"
	"starliner.app/internal/builder/infrastructure/git"
	buildergrpc "starliner.app/internal/builder/infrastructure/grpc"
	"starliner.app/internal/builder/infrastructure/nats/impl/queue"
	"starliner.app/internal/builder/presentation/grpc"
	builderqueue "starliner.app/internal/builder/presentation/queue"
	"starliner.app/internal/builder/presentation/scheduler"
	"starliner.app/internal/core/infrastructure/redis"
	"starliner.app/internal/core/infrastructure/s3"
	"starliner.app/internal/core/infrastructure/sentry"
)

func main() {
	fx.New(
		conf.Module,
		redis.Module,
		s3.Module,
		queue.Module,
		buildergrpc.Module,
		grpc.Module,
		git.Module,
		docker.Module,
		application.Module,
		builderqueue.Module,
		scheduler.Module,
		sentry.Module("builder"),
	).Run()
}
