package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "starliner.app/cmd/api/docs"
	"starliner.app/pkg/api/handler"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(
	rootHandler *handler.RootHandler) *Server {
	engine := gin.New()

	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/", rootHandler.GetRoot)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Server{engine: engine}
}

func (s *Server) Start() {
	s.engine.Run(":9090")
}
