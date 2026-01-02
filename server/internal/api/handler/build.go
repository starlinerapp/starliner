package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/service"
)

type BuildHandler struct {
	buildService *service.BuildService
}

func NewBuildHandler(buildService *service.BuildService) *BuildHandler {
	return &BuildHandler{
		buildService: buildService,
	}
}

// TriggerBuild FindAll godoc
// @Summary Trigger Build
// @Tags build
// @ID triggerBuild
// @Param X-User-ID header string true "User ID"
// @Product JSON
// @Success 200
// @Router /builds [post]
func (bh *BuildHandler) TriggerBuild(c *gin.Context) {
	err := bh.buildService.TriggerBuild()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusOK)
}
