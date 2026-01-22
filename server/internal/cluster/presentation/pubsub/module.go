package pubsub

import "go.uber.org/fx"

var Module = fx.Module(
	"scheduler",
	fx.Provide(NewWorker),
	fx.Invoke(RegisterWorker),
)
