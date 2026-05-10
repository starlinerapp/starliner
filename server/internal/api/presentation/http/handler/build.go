package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/response"
	"starliner.app/internal/api/presentation/http/sse"
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
// @State core
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

// StreamBuildLogs FindAll godoc
// @Summary Stream build logs
// @State core
// @Tags build
// @ID streamBuildLogs
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Build ID"
// @Product text/event-stream
// @Success 200
// @Header 200 {string} Content-Type "text/event-stream"
// @Header 200 {string} Cache-Control "no-cache"
// @Header 200 {string} Connection "keep-alive"
// @Router /builds/{id}/logs/stream [get]
func (bh *BuildHandler) StreamBuildLogs(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	buildId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	sw, ok := sse.NewWriter(c.Writer)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "streaming not supported"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	err = bh.buildApplication.StreamBuildLogs(c.Request.Context(), currentUser.Id, buildId, sw)
	if err != nil {
		sw.WriteError(err)
	}
}
