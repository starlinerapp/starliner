package service

import "go.uber.org/fx"

var Module = fx.Module(
	"service",
	fx.Provide(
		NewOrganizationService,
		NewEnvironmentService,
	),
)
