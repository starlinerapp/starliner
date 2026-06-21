package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/builder/application"
	"starliner.app/internal/builder/conf"
	"starliner.app/internal/builder/domain/service"
	buildergrpc "starliner.app/internal/builder/infrastructure/grpc"
	"starliner.app/internal/builder/infrastructure/nats/impl/queue"
	"starliner.app/internal/builder/presentation/grpc"
	builderqueue "starliner.app/internal/builder/presentation/queue"
	"starliner.app/internal/builder/presentation/scheduler"
	"starliner.app/internal/core/infrastructure/redis"
	"starliner.app/internal/core/infrastructure/sentry"
)

func main() {
	fx.New(
		conf.Module,
		redis.Module,
		queue.Module,
		buildergrpc.Module,
		grpc.Module,
		application.Module,
		service.Module,
		builderqueue.Module,
		scheduler.Module,
		sentry.Module("builder"),
	).Run()
}
