package application

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"application",
	fx.Provide(NewUserApplication),
	fx.Provide(NewEnvironmentApplication),
	fx.Provide(NewProjectApplication),
	fx.Provide(NewOrganizationApplication),
	fx.Provide(NewBuildApplication),
	fx.Provide(NewClusterApplication),
)
