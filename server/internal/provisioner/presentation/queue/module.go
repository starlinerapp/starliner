package queue

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"queue",
	fx.Provide(NewConsumer),
	fx.Invoke(RegisterConsumer),
)
