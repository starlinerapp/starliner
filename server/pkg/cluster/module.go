package cluster

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"cluster",
	fx.Provide(
		NewOrchestrator,
	),
	fx.Invoke(
		RegisterOrchestrator,
	),
)
