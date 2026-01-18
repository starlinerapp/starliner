package queue

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"provisioner",
	fx.Provide(NewConsumer),
	fx.Invoke(RegisterConsumer),
)
