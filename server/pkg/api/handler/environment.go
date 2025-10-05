package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/pkg/api/dto/request"
	"starliner.app/pkg/domain"
	"starliner.app/pkg/service"
)

type EnvironmentHandler struct {
	environmentService *service.EnvironmentService
}

func NewEnvironmentHandler(environmentService *service.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{environmentService: environmentService}
}

// CreateEnvironment FindAll godoc
// @Summary Create Environment
// @Tags environment
// @ID createEnvironment
// @Product JSON
// @Param data body request.CreateEnvironment true "Create Environment"
// @Success 201
// @Router /environments [post]
func (eh *EnvironmentHandler) CreateEnvironment(c *gin.Context) {
	currentUser := c.MustGet("user").(*domain.User)
	var env request.CreateEnvironment
	if err := c.BindJSON(&env); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_, err := eh.environmentService.CreateEnvironment(c, env.Name, currentUser.Id, env.OrganizationID, env.ProjectID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusCreated)
}
