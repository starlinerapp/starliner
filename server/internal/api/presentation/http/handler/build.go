package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/api/application"
)

type BuildHandler struct {
	buildApplication *application.BuildApplication
}

func NewBuildHandler(buildApplication *application.BuildApplication) *BuildHandler {
	return &BuildHandler{
		buildApplication: buildApplication,
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
	err := bh.buildApplication.TriggerBuild()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusOK)
}
