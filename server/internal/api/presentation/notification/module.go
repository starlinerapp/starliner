package notification

import "go.uber.org/fx"

var Module = fx.Module(
	"notification-consumer",
	fx.Provide(NewConsumer),
	fx.Invoke(RegisterConsumer),
)
