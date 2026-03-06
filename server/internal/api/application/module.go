package application

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"application",
	fx.Provide(
		NewUserApplication,
		NewEnvironmentApplication,
		NewProjectApplication,
		NewOrganizationApplication,
		NewClusterApplication,
		NewDeploymentApplication,
	),
)
