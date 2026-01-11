package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/application"
	"starliner.app/internal/domain/value"
	"starliner.app/internal/presentation/http/dto/request"
)

type DeploymentHandler struct {
	deploymentApplication *application.DeploymentApplication
}

func NewDeploymentHandler(
	deploymentApplication *application.DeploymentApplication,
) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentApplication: deploymentApplication,
	}
}

// DeployDatabase FindAll godoc
// @Summary Deploy databases
// @Tags deployment
// @ID deployDatabase
// @Param X-User-ID header string true "User ID"
// @Param data body request.DeployDatabase true "Deploy Database"
// @Product JSON
// @Success 200
// @Router /deployments/databases [post]
func (dh *DeploymentHandler) DeployDatabase(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	var body request.DeployDatabase
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := dh.deploymentApplication.DeployDatabase(
		c.Request.Context(),
		currentUser.Id,
		body.EnvironmentId,
		value.Database(body.Database),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}
