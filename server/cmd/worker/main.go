package main

import (
	"go.uber.org/fx"
	"starliner.app/internal/application"
	"starliner.app/internal/conf"
	"starliner.app/internal/domain/repository"
	"starliner.app/internal/domain/service"
	"starliner.app/internal/infrastructure/crypto"
	"starliner.app/internal/infrastructure/dagger"
	"starliner.app/internal/infrastructure/nats/queue"
	"starliner.app/internal/infrastructure/postgres"
	"starliner.app/internal/infrastructure/s3"
	"starliner.app/internal/infrastructure/ssh"
	"starliner.app/internal/presentation/queue/builder"
	"starliner.app/internal/presentation/queue/cluster"
	"starliner.app/internal/presentation/queue/provisioner"
)

func main() {
	fx.New(
		conf.Module,
		postgres.Module,
		queue.Module,
		s3.Module,
		dagger.Module,
		crypto.Module,
		ssh.Module,
		repository.Module,
		application.Module,
		service.Module,
		builder.Module,
		provisioner.Module,
		cluster.Module,
	).Run()
}
