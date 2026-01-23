package health

import "go.uber.org/fx"

var Module = fx.Module(
	"pubsub",
	fx.Provide(NewSubscriber),
	fx.Invoke(RegisterSubscriber),
)
