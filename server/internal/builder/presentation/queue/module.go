package queue

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"builder",
	fx.Provide(
		NewConsumer,
	),
	fx.Invoke(
		RegisterConsumer,
	),
)
