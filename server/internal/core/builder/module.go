package builder

import (
	"go.uber.org/fx"
	"starliner.app/internal/infrastructure/dagger"
	"starliner.app/internal/infrastructure/objectstore"
)

var Module = fx.Module(
	"build",
	objectstore.Module,
	dagger.Module,
	fx.Provide(
		NewOrchestrator,
	),
	fx.Invoke(
		RegisterOrchestrator,
	),
)
