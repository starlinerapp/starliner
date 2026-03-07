package k8s

import "go.uber.org/fx"

var Module = fx.Module(
	"k8s",
	fx.Provide(
		NewHealth,
		NewSecret,
		NewLogs,
	),
)
