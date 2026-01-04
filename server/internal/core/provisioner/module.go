package provisioner

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"provisioner",
	fx.Provide(NewOrchestrator),
	fx.Invoke(RegisterOrchestrator),
)
