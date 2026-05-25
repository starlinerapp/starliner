package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
	"starliner.app/internal/api/presentation/http/dto/response"
	"starliner.app/internal/api/presentation/http/sse"
)

type EnvironmentHandler struct {
	environmentApplication *application.EnvironmentApplication
	notificationHub        *sse.EnvironmentNotificationHub
}

func NewEnvironmentHandler(environmentApplication *application.EnvironmentApplication, notificationHub *sse.EnvironmentNotificationHub) *EnvironmentHandler {
	return &EnvironmentHandler{environmentApplication: environmentApplication, notificationHub: notificationHub}
}

// CreateEnvironment FindAll godoc
// @Summary Create Environment
// @State core
// @Tags environment
// @ID createEnvironment
// @Param X-User-ID header string true "User ID"
// @Product JSON
// @Param data body request.CreateEnvironment true "Create Environment"
// @Success 201 {object} response.Environment
// @Router /environments [post]
func (eh *EnvironmentHandler) CreateEnvironment(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	var env request.CreateEnvironment
	if err := c.BindJSON(&env); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newEnv, err := eh.environmentApplication.CreateEnvironment(c.Request.Context(), env.Name, currentUser.Id, env.OrganizationID, env.ProjectID, env.SourceEnvironmentID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusCreated, response.NewEnvironment(newEnv))
}

// DeleteEnvironment FindAll godoc
// @Summary Delete Environment
// @State core
// @Tags environment
// @ID deleteEnvironment
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Environment ID"
// @Success 200
// @Router /environments/{id} [delete]
func (eh *EnvironmentHandler) DeleteEnvironment(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	environmentId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err = eh.environmentApplication.DeleteEnvironment(c.Request.Context(), currentUser.Id, environmentId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusOK)
}

// GetEnvironmentDeployments FindAll godoc
// @Summary Get Environment Deployments
// @State core
// @Tags environment
// @ID getEnvironmentDeployments
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Environment ID"
// @Success 200 {object} response.Deployments
// @Router /environments/{id}/deployments [get]
func (eh *EnvironmentHandler) GetEnvironmentDeployments(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	environmentId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	deployments, err := eh.environmentApplication.GetEnvironmentDeployments(c.Request.Context(), environmentId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.NewDeployments(deployments))
}

// GetEnvironmentBuilds FindAll godoc
// @Summary Get Environment Builds
// @State core
// @Tags environment
// @ID getEnvironmentBuilds
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Environment ID"
// @Success 200 {array} response.GitDeploymentBuild
// @Router /environments/{id}/builds [get]
func (eh *EnvironmentHandler) GetEnvironmentBuilds(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	environmentId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	builds, err := eh.environmentApplication.GetEnvironmentGitDeploymentBuilds(c.Request.Context(), currentUser.Id, environmentId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.NewGitDeploymentBuilds(builds))
}

// GetEnvironmentConnectedBranch FindAll godoc
// @Summary Get Environment Connected Branch
// @State core
// @Tags environment
// @ID getEnvironmentConnectedBranch
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Environment ID"
// @Success 200 {object} response.EnvironmentBranch
// @Router /environments/{id}/branch [get]
func (eh *EnvironmentHandler) GetEnvironmentConnectedBranch(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	environmentId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	branch, err := eh.environmentApplication.GetEnvironmentBranch(c.Request.Context(), currentUser.Id, environmentId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.EnvironmentBranch{Branch: branch})
}

// UpdateEnvironmentConnectedBranch FindAll godoc
// @Summary Update Environment Connected Branch
// @State core
// @Tags environment
// @ID updateEnvironmentConnectedBranch
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Environment ID"
// @Product JSON
// @Param data body request.UpdateEnvironmentConnectBranch true "Update Environment Connected Branch"
// @Success 200
// @Router /environments/{id}/branch [put]
func (eh *EnvironmentHandler) UpdateEnvironmentConnectedBranch(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	environmentId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	var branch request.UpdateEnvironmentConnectBranch
	if err := c.BindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = eh.environmentApplication.UpdateEnvironmentBranch(c.Request.Context(), currentUser.Id, environmentId, branch.Branch)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Status(http.StatusOK)
}

// StreamEnvironmentNotifications godoc
// @Summary Stream environment notifications
// @State core
// @Tags environment
// @ID streamEnvironmentNotifications
// @Param X-User-ID header string true "User ID"
// @Param X-Correlation-ID header string true "Correlation ID"
// @Param id path int true "Environment ID"
// @Product text/event-stream
// @Success 200
// @Header 200 {string} Content-Type "text/event-stream"
// @Header 200 {string} Cache-Control "no-cache"
// @Header 200 {string} Connection "keep-alive"
// @Router /environments/{id}/notifications [get]
func (eh *EnvironmentHandler) StreamEnvironmentNotifications(c *gin.Context) {
	correlationId := c.GetHeader("X-Correlation-ID")
	if correlationId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing X-Correlation-ID header"})
		return
	}

	environmentId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid environment id"})
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

	ch := eh.notificationHub.Subscribe(correlationId, environmentId)
	defer eh.notificationHub.Unsubscribe(correlationId, environmentId, ch)

	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case notification := <-ch:
			eh.notificationHub.WriteNotification(sw, notification)
		}
	}
}
