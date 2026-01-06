package http

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"log"
	"net/http"
	_ "starliner.app/cmd/api/docs"
	handler2 "starliner.app/internal/presentation/http/handler"
	middleware2 "starliner.app/internal/presentation/http/middleware"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(
	auth *middleware2.BasicAuthMiddleware,
	user *middleware2.UserMiddleware,
	rootHandler *handler2.RootHandler,
	userHandler *handler2.UserHandler,
	organizationHandler *handler2.OrganizationHandler,
	projectHandler *handler2.ProjectHandler,
	environmentHandler *handler2.EnvironmentHandler,
	buildHandler *handler2.BuildHandler,
	clusterHandler *handler2.ClusterHandler,
) *Server {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	engine.Use(auth.WithBasicAuth(), user.WithUser())
	engine.GET("/", rootHandler.GetRoot)
	engine.GET("/me", userHandler.GetUser)

	organizationRoutes := engine.Group("/organizations")
	{
		organizationRoutes.POST("", organizationHandler.CreateOrganization)
		organizationRoutes.GET("", organizationHandler.GetUserOrganizations)
		organizationRoutes.GET("/:id/projects", organizationHandler.GetOrganizationProjects)
		organizationRoutes.GET("/:id/clusters", organizationHandler.GetOrganizationClusters)
	}

	projectRoutes := engine.Group("/projects")
	{
		projectRoutes.POST("", projectHandler.CreateProject)
		projectRoutes.GET("/:id", projectHandler.GetProject)
		projectRoutes.DELETE("/:id", projectHandler.DeleteProject)
	}

	environmentRoutes := engine.Group("/environments")
	{
		environmentRoutes.POST("", environmentHandler.CreateEnvironment)
	}

	buildRoutes := engine.Group("/builds")
	{
		buildRoutes.POST("", buildHandler.TriggerBuild)
	}

	clusterRoutes := engine.Group("/clusters")
	{
		clusterRoutes.POST("", clusterHandler.CreateCluster)
		clusterRoutes.GET("/:id", clusterHandler.GetCluster)
		clusterRoutes.GET("/:id/private-key", clusterHandler.GetClusterPrivateKey)
		clusterRoutes.DELETE("/:id", clusterHandler.DeleteCluster)
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
