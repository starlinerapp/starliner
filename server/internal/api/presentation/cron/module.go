package cron

import "go.uber.org/fx"

var Module = fx.Module(
	"cron",
	fx.Provide(
		NewScheduler,
		NewCron,
	),
	fx.Invoke(RegisterCron),
)
