package cluster

import (
	"go.uber.org/fx"
	"starliner.app/internal/infrastructure/db"
	"starliner.app/internal/repository"
)

var Module = fx.Module(
	"cluster",
	db.Module,
	repository.Module,
	fx.Provide(
		NewOrchestrator,
	),
	fx.Invoke(
		RegisterOrchestrator,
	),
)
