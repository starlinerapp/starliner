package queue

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"cluster",
	fx.Provide(NewConsumer),
	fx.Invoke(RegisterConsumer),
)
