package builder

import (
	"go.uber.org/fx"
	"starliner.app/pkg/dagger"
	"starliner.app/pkg/objectstore"
	"starliner.app/pkg/queue"
)

var Module = fx.Module(
	"build",
	objectstore.Module,
	queue.Module,
	dagger.Module,
	fx.Provide(
		NewOrchestrator,
	),
	fx.Invoke(
		RegisterOrchestrator,
	),
)
