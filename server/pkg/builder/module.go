package builder

import (
	"go.uber.org/fx"
	"starliner.app/pkg/objectstore"
	"starliner.app/pkg/queue"
)

var Module = fx.Module(
	"build",
	objectstore.Module,
	queue.Module,
	fx.Provide(
		NewOrchestrator,
	),
	fx.Invoke(
		RegisterOrchestrator,
	),
)
