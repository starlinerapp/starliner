package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
	"starliner.app/internal/api/presentation/http/dto/response"
	"strconv"
)

type EnvironmentHandler struct {
	environmentApplication *application.EnvironmentApplication
}

func NewEnvironmentHandler(environmentApplication *application.EnvironmentApplication) *EnvironmentHandler {
	return &EnvironmentHandler{environmentApplication: environmentApplication}
}

// CreateEnvironment FindAll godoc
// @Summary Create Environment
// @Tags environment
// @ID createEnvironment
// @Param X-User-ID header string true "User ID"
// @Product JSON
// @Param data body request.CreateEnvironment true "Create Environment"
// @Success 201
// @Router /environments [post]
func (eh *EnvironmentHandler) CreateEnvironment(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	var env request.CreateEnvironment
	if err := c.BindJSON(&env); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := eh.environmentApplication.CreateEnvironment(c.Request.Context(), env.Name, currentUser.Id, env.OrganizationID, env.ProjectID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusCreated)
}

// GetEnvironmentDeployments FindAll godoc
// @Summary Get Environment Deployments
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
