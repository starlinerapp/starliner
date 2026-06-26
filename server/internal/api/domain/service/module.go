package service

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"service",
	fx.Provide(
		NewBuildService,
		NewOrganizationService,
		NewEnvironmentService,
		NewDeploymentService,
		NewGitDeploymentService,
		NewTeamService,
		NewClusterService,
		NewParserService,
		NewResolverService,
		NewTokenService,
	),
)
