package handler

import "go.uber.org/fx"

var Module = fx.Module(
	"handler",
	fx.Provide(
		NewRootHandler,
		NewUserHandler,
		NewEnvironmentHandler,
		NewProjectHandler,
		NewOrganizationHandler,
		NewClusterHandler,
		NewDeploymentHandler,
		NewBuildHandler,
		NewTeamHandler,
		NewGithubHandler,
		NewGithubAppHandler,
		NewWebhookHandler,
		NewInternalHandler,
	),
)
