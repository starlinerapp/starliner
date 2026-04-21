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

type ProjectHandler struct {
	projectApplication *application.ProjectApplication
}

func NewProjectHandler(projectApplication *application.ProjectApplication) *ProjectHandler {
	return &ProjectHandler{
		projectApplication: projectApplication,
	}
}

// CreateProject FindAll godoc
// @Summary Create Project
// @Tags project
// @ID createProject
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param data body request.CreateProject true "Create Project"
// @Success 200 {object} response.Project
// @Router /projects [post]
func (ph *ProjectHandler) CreateProject(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	var project request.CreateProject
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newProject, err := ph.projectApplication.CreateProject(c.Request.Context(), project.Name, project.OrganizationId, project.ClusterId, currentUser.Id, project.TeamId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusCreated, response.NewProject(newProject))
}

// GetProject FindAll godoc
// @Summary Get Project
// @Tags project
// @ID getProject
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Project ID"
// @Success 200 {object} response.Project
// @Router /projects/{id} [get]
func (ph *ProjectHandler) GetProject(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	projectId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	project, err := ph.projectApplication.GetProject(c.Request.Context(), projectId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.NewProject(project))
}

// DeleteProject FindAll godoc
// @Summary Delete Project
// @Tags project
// @ID deleteProject
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Project ID"
// @Success 200
// @Router /projects/{id} [delete]
func (ph *ProjectHandler) DeleteProject(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	projectId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = ph.projectApplication.DeleteProject(c.Request.Context(), projectId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	c.Status(http.StatusOK)
}

// GetProjectCluster FindAll godoc
// @Summary Get Project Cluster
// @Tags project
// @ID getProjectCluster
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Project ID"
// @Success 200 {object} response.ProjectCluster
// @Router /projects/{id}/cluster [get]
func (ph *ProjectHandler) GetProjectCluster(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	projectId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	cluster, err := ph.projectApplication.GetProjectCluster(c.Request.Context(), projectId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	c.JSON(http.StatusOK, response.NewProjectCluster(cluster))
}

// GetProjectEnvironments godoc
// @Summary Get Project Environments
// @Tags project
// @ID getProjectEnvironments
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Project ID"
// @Success 200 {array} response.Environment
// @Router /projects/{id}/environments [get]
func (ph *ProjectHandler) GetProjectEnvironments(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	projectId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	environments, err := ph.projectApplication.GetProjectEnvironments(c.Request.Context(), projectId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.NewEnvironments(environments))
}

// GetProjectPreviewEnvironmentEnabled godoc
// @Summary Get Project Preview Environment Enabled
// @Tags project
// @ID getProjectPreviewEnvironmentEnabled
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Project ID"
// @Success 200 {object} response.ProjectPreviewEnvironmentEnabled
// @Router /projects/{id}/preview-environment/enabled [get]
func (ph *ProjectHandler) GetProjectPreviewEnvironmentEnabled(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	projectId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	enabled, err := ph.projectApplication.GetProjectPreviewEnvironmentEnabled(c.Request.Context(), projectId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.ProjectPreviewEnvironmentEnabled{Enabled: enabled})
}

// ToggleProjectPreviewEnvironmentEnabled godoc
// @Summary Toggle Project Preview Environment Enabled
// @Tags project
// @ID toggleProjectPreviewEnvironmentEnabled
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Project ID"
// @Success 200 {object} response.ProjectPreviewEnvironmentEnabled
// @Router /projects/{id}/preview-environment/enabled [put]
func (ph *ProjectHandler) ToggleProjectPreviewEnvironmentEnabled(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	projectId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	enabled, err := ph.projectApplication.ToggleProjectPreviewEnvironmentEnabled(c.Request.Context(), projectId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.ProjectPreviewEnvironmentEnabled{Enabled: enabled})
}
