package queue

import "go.uber.org/fx"

var Module = fx.Module(
	"runner-queue",
	fx.Provide(NewConsumer),
	fx.Invoke(RegisterConsumer),
)
