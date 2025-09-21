package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "starliner.app/cmd/api/docs"
	"starliner.app/pkg/api/handler"
	"starliner.app/pkg/api/middleware"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(
	rootHandler *handler.RootHandler,
	auth *middleware.BasicAuthMiddleware,
) *Server {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	engine.Use(auth.BasicAuth())
	engine.GET("/", rootHandler.GetRoot)

	return &Server{engine: engine}
}

func (s *Server) Start() {
	s.engine.Run(":9090")
}
