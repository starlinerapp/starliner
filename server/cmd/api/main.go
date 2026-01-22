package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/repository"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/infrastructure/nats/impl/pubsub"
	"starliner.app/internal/api/infrastructure/nats/impl/queue"
	"starliner.app/internal/api/infrastructure/postgres"
	"starliner.app/internal/api/presentation/http"
	clusterqueue "starliner.app/internal/api/presentation/queue/cluster"
	"starliner.app/internal/api/presentation/scheduler"
	"starliner.app/internal/core/infrastructure/crypto"
	"starliner.app/internal/core/infrastructure/s3"
	"starliner.app/internal/provisioner/infrastructure/ssh"
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
		ssh.Module,
		repository.Module,
		application.Module,
		service.Module,
		http.Module,
		clusterqueue.Module,
		scheduler.Module,
	).Run()
}
