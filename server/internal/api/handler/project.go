package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/api/dto/request"
	"starliner.app/internal/api/dto/response"
	"starliner.app/internal/service"
	"starliner.app/internal/service/model"
	"strconv"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
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
	currentUser := c.MustGet("user").(*model.User)
	var project request.CreateProject
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newProject, err := ph.projectService.CreateProject(c.Request.Context(), project.Name, project.OrganizationId, project.ClusterId, currentUser.Id)
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
	currentUser := c.MustGet("user").(*model.User)
	projectId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	project, err := ph.projectService.GetProject(c.Request.Context(), projectId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.NewProject(project))
}
