package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/application"
	"starliner.app/internal/domain/value"
	"starliner.app/internal/presentation/http/dto/request"
)

type EnvironmentHandler struct {
	environmentService *application.EnvironmentApplication
}

func NewEnvironmentHandler(environmentService *application.EnvironmentApplication) *EnvironmentHandler {
	return &EnvironmentHandler{environmentService: environmentService}
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

	err := eh.environmentService.CreateEnvironment(c.Request.Context(), env.Name, currentUser.Id, env.OrganizationID, env.ProjectID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusCreated)
}
