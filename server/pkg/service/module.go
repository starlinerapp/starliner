package service

import "go.uber.org/fx"

var Module = fx.Module(
	"service",
	fx.Provide(NewUserService),
	fx.Provide(NewEnvironmentService),
	fx.Provide(NewProjectService),
	fx.Provide(NewOrganizationService),
)
