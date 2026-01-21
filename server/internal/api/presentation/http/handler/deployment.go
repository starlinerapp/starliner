package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
	coreValue "starliner.app/internal/core/domain/value"
	"strconv"
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
// @Summary Deploy database
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
		coreValue.Database(body.Database),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}

// DeleteDatabase FindAll godoc
// @Summary Delete database
// @Tags deployment
// @ID deleteDatabase
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Deployment ID"
// @Product JSON
// @Success 200
// @Router /deployments/databases/{id} [delete]
func (dh *DeploymentHandler) DeleteDatabase(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	deploymentId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = dh.deploymentApplication.DeleteDatabase(
		c.Request.Context(),
		deploymentId,
		currentUser.Id,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}
