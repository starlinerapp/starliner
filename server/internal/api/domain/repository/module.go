package repository

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"repository",
	fx.Provide(
		NewEnvironmentRepository,
		NewOrganizationRepository,
		NewProjectRepository,
		NewUserRepository,
		NewDeploymentRepository,
		NewClusterRepository,
		NewBuildRepository,
		NewTeamRepository,
		NewGithubAppRepository,
	),
)
