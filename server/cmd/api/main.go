package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/repository"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/infrastructure/grpc"
	"starliner.app/internal/api/infrastructure/nats/impl/pubsub"
	"starliner.app/internal/api/infrastructure/nats/impl/queue"
	"starliner.app/internal/api/infrastructure/postgres"
	"starliner.app/internal/api/presentation/http"
	clusterpubsub "starliner.app/internal/api/presentation/pubsub/health"
	clusterqueue "starliner.app/internal/api/presentation/queue/cluster"
	"starliner.app/internal/api/presentation/scheduler"
	coreService "starliner.app/internal/core/domain/service"
	"starliner.app/internal/core/infrastructure/crypto"
	"starliner.app/internal/core/infrastructure/s3"
)

// @title Starliner API
// @version 1.0
// @securityDefinitions.basic BasicAuth
func main() {
	fx.New(
		conf.Module,
		postgres.Module,
		queue.Module,
		pubsub.Module,
		s3.Module,
		crypto.Module,
		grpc.Module,
		repository.Module,
		coreService.Module,
		service.Module,
		application.Module,
		http.Module,
		clusterqueue.Module,
		clusterpubsub.Module,
		scheduler.Module,
	).Run()
}
