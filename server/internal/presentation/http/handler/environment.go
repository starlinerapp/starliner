package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/application"
	"starliner.app/internal/domain/value"
	"starliner.app/internal/presentation/http/dto/request"
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

// DeployDatabase FindAll godoc
// @Summary Deploy database to environment
// @Tags environment
// @ID deployDatabase
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Environment ID"
// @Param data body request.DeployDatabase true "Deploy Database"
// @Product JSON
// @Success 200
// @Router /environments/{id}/databases [post]
func (eh *EnvironmentHandler) DeployDatabase(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	environmentId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	var body request.DeployDatabase
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = eh.environmentApplication.DeployDatabase(
		c.Request.Context(),
		currentUser.Id,
		environmentId,
		value.Database(body.Database),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}
