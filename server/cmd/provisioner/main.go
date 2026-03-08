package main

import (
	"go.uber.org/fx"
	coreService "starliner.app/internal/core/domain/service"
	"starliner.app/internal/core/infrastructure/crypto"
	"starliner.app/internal/provisioner/application"
	"starliner.app/internal/provisioner/conf"
	"starliner.app/internal/provisioner/infrastructure/ansible"
	"starliner.app/internal/provisioner/infrastructure/nats/impl/queue"
	"starliner.app/internal/provisioner/infrastructure/pulumi"
	"starliner.app/internal/provisioner/infrastructure/ssh"
	provisionerqueue "starliner.app/internal/provisioner/presentation/queue"
)

func main() {
	fx.New(
		conf.Module,
		crypto.Module,
		ssh.Module,
		queue.Module,
		pulumi.Module,
		ansible.Module,
		coreService.Module,
		application.Module,
		provisionerqueue.Module,
	).Run()
}
