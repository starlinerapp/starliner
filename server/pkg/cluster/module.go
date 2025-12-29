package cluster

import (
	"go.uber.org/fx"
	"starliner.app/pkg/db"
	"starliner.app/pkg/repository"
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
