package application

import (
	"go.uber.org/fx"
	"starliner.app/internal/cluster/domain/port"
)

var Module = fx.Module(
	"application",
	fx.Provide(
		NewDeploymentApplication,
		NewImageApplication,
		NewDatabaseApplication,
		NewIngressApplication,
		NewStatusApplication,
		NewLogsApplication,
		NewDeploymentStatusApplication,
		func(da *DeploymentStatusApplication) port.LogPublisher { return da },
		NewTTYApplication,
	),
)
