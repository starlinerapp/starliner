package scheduler

import "go.uber.org/fx"

var Module = fx.Module(
	"scheduler",
	fx.Provide(NewScheduler),
	fx.Invoke(RegisterScheduler),
)
