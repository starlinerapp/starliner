package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/response"
)

type BuildHandler struct {
	buildApplication *application.BuildApplication
}

func NewBuildHandler(buildApplication *application.BuildApplication) *BuildHandler {
	return &BuildHandler{
		buildApplication: buildApplication,
	}
}

// GetBuildLogs FindAll godoc
// @Summary Get Build Logs
// @Tags build
// @ID getBuildLogs
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Build ID"
// @Success 200 {object} response.BuildLogs
// @Router /builds/{id}/logs [get]
func (bh *BuildHandler) GetBuildLogs(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	buildId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	logs, err := bh.buildApplication.GetBuildLogs(c.Request.Context(), currentUser.Id, buildId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, response.NewBuildLogs(logs))
}
