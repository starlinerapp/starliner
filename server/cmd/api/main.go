package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/application"
	"starliner.app/internal/conf"
	"starliner.app/internal/domain/repository"
	"starliner.app/internal/domain/service"
	"starliner.app/internal/infrastructure/ansible"
	"starliner.app/internal/infrastructure/crypto"
	"starliner.app/internal/infrastructure/dagger"
	"starliner.app/internal/infrastructure/dagger/impl/docker"
	"starliner.app/internal/infrastructure/nats/impl/queue"
	"starliner.app/internal/infrastructure/postgres"
	"starliner.app/internal/infrastructure/s3"
	"starliner.app/internal/infrastructure/ssh"
	"starliner.app/internal/presentation/http"
)

// @title Starliner API
// @version 1.0
// @securityDefinitions.basic BasicAuth
func main() {
	fx.New(
		conf.Module,
		postgres.Module,
		dagger.Module,
		docker.Module,
		ansible.Module,
		queue.Module,
		s3.Module,
		crypto.Module,
		ssh.Module,
		repository.Module,
		application.Module,
		service.Module,
		http.Module,
	).Run()
}
