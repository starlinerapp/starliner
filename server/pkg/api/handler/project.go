package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/pkg/api/dto/request"
	"starliner.app/pkg/domain"
	"starliner.app/pkg/service"
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
// @Param data body request.CreateProject true "Create Project"
// @Success 200
// @Router /projects [post]
func (ph *ProjectHandler) CreateProject(c *gin.Context) {
	currentUser := c.MustGet("user").(*domain.User)
	var project request.CreateProject
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_, err := ph.projectService.CreateProject(c, project.Name, project.OrganizationId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusCreated)
}
