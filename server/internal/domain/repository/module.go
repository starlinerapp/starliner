package repository

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"repository",
	fx.Provide(NewEnvironmentRepository),
	fx.Provide(NewOrganizationRepository),
	fx.Provide(NewProjectRepository),
	fx.Provide(NewUserRepository),
	fx.Provide(NewClusterRepository),
)
