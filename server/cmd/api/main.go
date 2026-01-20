package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/repository"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/infrastructure/postgres"
	"starliner.app/internal/api/presentation/http"
	"starliner.app/internal/api/presentation/queue/cluster"
	"starliner.app/internal/core/infrastructure/crypto"
	nats "starliner.app/internal/core/infrastructure/nats/impl/queue"
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
		nats.Module,
		s3.Module,
		crypto.Module,
		ssh.Module,
		repository.Module,
		application.Module,
		service.Module,
		http.Module,
		queue.Module,
	).Run()
}
