package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	_ "starliner.app/cmd/api/docs"
	"starliner.app/internal/api/presentation/http/handler"
	"starliner.app/internal/api/presentation/http/middleware"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(
	auth *middleware.BasicAuthMiddleware,
	user *middleware.UserMiddleware,
	rootHandler *handler.RootHandler,
	userHandler *handler.UserHandler,
	organizationHandler *handler.OrganizationHandler,
	projectHandler *handler.ProjectHandler,
	environmentHandler *handler.EnvironmentHandler,
	clusterHandler *handler.ClusterHandler,
	deploymentHandler *handler.DeploymentHandler,
	buildHandler *handler.BuildHandler,
	teamHandler *handler.TeamHandler,
	githubHandler *handler.GithubHandler,
	githubAppHandler *handler.GithubAppHandler,
	webhookHandler *handler.WebhookHandler,
) *Server {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	webhookRoutes := engine.Group("/webhooks")
	{
		webhookRoutes.POST("/github", webhookHandler.HandleGithubWebhook)
	}

	engine.Use(auth.WithBasicAuth(), user.WithUser())
	engine.GET("/", rootHandler.GetRoot)
	engine.GET("/me", userHandler.GetUser)

	organizationRoutes := engine.Group("/organizations")
	{
		organizationRoutes.POST("", organizationHandler.CreateOrganization)
		organizationRoutes.GET("", organizationHandler.GetUserOrganizations)
		organizationRoutes.GET("/:id/projects", organizationHandler.GetUserProjects)
		organizationRoutes.GET("/:id/clusters", organizationHandler.GetOrganizationClusters)
		organizationRoutes.POST("/:id/settings/credential/hetzner", organizationHandler.UpsertHetznerCredential)
		organizationRoutes.GET("/:id/settings/credential/hetzner", organizationHandler.GetHetznerCredential)
		organizationRoutes.POST("/:id/invites", organizationHandler.SendEmailInvite)
		organizationRoutes.POST("/:id/teams", teamHandler.CreateTeam)
		organizationRoutes.GET("/:id/teams", teamHandler.GetUserTeams)
		organizationRoutes.POST("/:id/teams/join", teamHandler.JoinTeam)
		organizationRoutes.GET("/:id/members", organizationHandler.GetOrganizationMembers)
	}

	inviteRoutes := engine.Group("/invites")
	{
		inviteRoutes.GET("/:inviteId", organizationHandler.GetInviteDetails)
		inviteRoutes.POST("/accept", organizationHandler.AcceptInvite)
	}

	projectRoutes := engine.Group("/projects")
	{
		projectRoutes.POST("", projectHandler.CreateProject)
		projectRoutes.GET("/:id", projectHandler.GetProject)
		projectRoutes.DELETE("/:id", projectHandler.DeleteProject)
		projectRoutes.GET("/:id/cluster", projectHandler.GetProjectCluster)
		projectRoutes.GET("/:id/environments", projectHandler.GetProjectEnvironments)
		projectRoutes.GET("/:id/preview-environment/enabled", projectHandler.GetProjectPreviewEnvironmentEnabled)
		projectRoutes.PUT("/:id/preview-environment/enabled", projectHandler.ToggleProjectPreviewEnvironmentEnabled)
	}

	environmentRoutes := engine.Group("/environments")
	{
		environmentRoutes.POST("", environmentHandler.CreateEnvironment)
		environmentRoutes.GET("/:id/deployments", environmentHandler.GetEnvironmentDeployments)
		environmentRoutes.GET("/:id/builds", environmentHandler.GetEnvironmentBuilds)
		environmentRoutes.GET("/:id/branch", environmentHandler.GetEnvironmentConnectedBranch)
		environmentRoutes.PUT("/:id/branch", environmentHandler.UpdateEnvironmentConnectedBranch)
	}

	clusterRoutes := engine.Group("/clusters")
	{
		clusterRoutes.POST("", clusterHandler.CreateCluster)
		clusterRoutes.GET("/:id", clusterHandler.GetCluster)
		clusterRoutes.GET("/:id/private-key", clusterHandler.GetClusterPrivateKey)
		clusterRoutes.DELETE("/:id", clusterHandler.DeleteCluster)
	}

	deploymentRoutes := engine.Group("/deployments")
	{
		deploymentRoutes.POST("/git", deploymentHandler.DeployFromGitRepository)
		deploymentRoutes.PUT("/git/:deploymentId", deploymentHandler.UpdateDeployFromGitRepository)
		deploymentRoutes.POST("/images", deploymentHandler.DeployImage)
		deploymentRoutes.PUT("/images/:deploymentId", deploymentHandler.UpdateImageDeployment)
		deploymentRoutes.POST("/databases", deploymentHandler.DeployDatabase)
		deploymentRoutes.POST("/ingresses", deploymentHandler.DeployIngress)
		deploymentRoutes.PUT("/ingresses/:deploymentId", deploymentHandler.UpdateIngressDeployment)
		deploymentRoutes.DELETE("/:id", deploymentHandler.DeleteDeployment)
		deploymentRoutes.GET("/:id/logs", deploymentHandler.StreamDeploymentLogs)
	}

	buildRoutes := engine.Group("/builds")
	{
		buildRoutes.GET("/:id/logs", buildHandler.GetBuildLogs)
		buildRoutes.GET("/:id/logs/stream", buildHandler.StreamBuildLogs)
	}

	webSocketRoutes := engine.Group("/ws")
	{
		webSocketRoutes.GET("/deployments/:id", deploymentHandler.OpenTTY)
		webSocketRoutes.GET("/clusters/:id", clusterHandler.OpenTTY)
	}

	teamRoutes := engine.Group("/teams")
	{
		teamRoutes.GET("/:teamId/members", teamHandler.GetTeamMembers)
		teamRoutes.POST("/:teamId/members", teamHandler.AddTeamMember)
		teamRoutes.DELETE("/:teamId/members", teamHandler.RemoveTeamMember)
		teamRoutes.GET("/:teamId/repos", teamHandler.GetTeamRepositories)
		teamRoutes.POST("/:teamId/repos", teamHandler.AssignRepoToTeam)
		teamRoutes.DELETE("/:teamId/repos/:repoId", teamHandler.UnassignRepoFromTeam)
		teamRoutes.GET("/:teamId/clusters", teamHandler.GetTeamClusters)
		teamRoutes.POST("/:teamId/clusters/:clusterId", teamHandler.AssignClusterToTeam)
		teamRoutes.DELETE("/:teamId/clusters/:clusterId", teamHandler.UnassignClusterFromTeam)
	}

	githubRoutes := engine.Group("/github")
	{
		githubRoutes.GET("/repositories/:organizationId", githubHandler.GetRepositories)
		githubRoutes.GET("/all-repositories/:organizationId", githubHandler.GetAllRepositories)
		githubRoutes.GET("/repositories/:organizationId/:owner/:repository/contents", githubHandler.GetRepositoryContents)
		githubRoutes.GET("/repositories/:organizationId/:owner/:repository/file", githubHandler.GetFileContent)
	}

	githubAppRoutes := engine.Group("/githubapps")
	{
		githubAppRoutes.POST("", githubAppHandler.CreateGithubApp)
		githubAppRoutes.GET("/:organizationId", githubAppHandler.GetGithubApp)
	}

	return &Server{engine: engine}
}

func RegisterServer(lc fx.Lifecycle, s *Server) {
	server := &http.Server{
		Addr:    ":9090",
		Handler: s.engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("failed to start server: %v", err)
				}
			}()
			log.Printf("Server listening on port 9090")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("Shutting down server...")
			return server.Shutdown(ctx)
		},
	})
}
