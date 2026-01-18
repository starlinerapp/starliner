package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/repository"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/presentation/http"
	"starliner.app/internal/core/conf"
	coreRepository "starliner.app/internal/core/domain/repository"
	"starliner.app/internal/core/infrastructure/crypto"
	"starliner.app/internal/core/infrastructure/nats/impl/queue"
	"starliner.app/internal/core/infrastructure/postgres"
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
		s3.Module,
		crypto.Module,
		ssh.Module,
		coreRepository.Module,
		repository.Module,
		application.Module,
		service.Module,
		http.Module,
	).Run()
}
