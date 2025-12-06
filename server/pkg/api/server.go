package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	_ "starliner.app/cmd/api/docs"
	"starliner.app/pkg/api/handler"
	"starliner.app/pkg/api/middleware"
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
	}

	projectRoutes := engine.Group("/projects")
	{
		projectRoutes.POST("", projectHandler.CreateProject)
		projectRoutes.GET("", projectHandler.GetProject)
	}

	environmentRoutes := engine.Group("/environments")
	{
		environmentRoutes.POST("", environmentHandler.CreateEnvironment)
	}

	return &Server{engine: engine}
}

func (s *Server) Start() {
	err := s.engine.Run(":9090")
	if err != nil {
		log.Fatalln("Failed to start server")
	}
}
