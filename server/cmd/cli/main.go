package main

import (
	"fmt"
	"os"

	"go.uber.org/fx"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/repository"
	"starliner.app/internal/api/domain/service"
	"starliner.app/internal/api/infrastructure/nats/impl/queue"
	"starliner.app/internal/api/infrastructure/postgres"
	"starliner.app/internal/api/presentation/cli"
)

func main() {
	app := fx.New(
		fx.NopLogger,
		conf.Module,
		postgres.Module,
		queue.Module,
		repository.Module,
		service.Module,
		application.Module,
		cli.Module,
	)
	if err := app.Err(); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			return
		}
		os.Exit(1)
	}
	app.Run()
}
