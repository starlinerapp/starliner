package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
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

// DeployImage FindAll godoc
// @Summary Deploy image
// @Tags deployment
// @ID deployImage
// @Param X-User-ID header string true "User ID"
// @Param data body request.DeployImage true "Deploy Image"
// @Product JSON
// @Success 200
// @Router /deployments/images [post]
func (dh *DeploymentHandler) DeployImage(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	var body request.DeployImage
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := dh.deploymentApplication.DeployImage(
		c.Request.Context(),
		currentUser.Id,
		body.EnvironmentId,
		body.ServiceName,
		body.ImageName,
		body.Tag,
		body.Port,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
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
		value.Database(body.Database),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}

// DeleteDeployment FindAll godoc
// @Summary Delete deployment
// @Tags deployment
// @ID deleteDeployment
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Deployment ID"
// @Product JSON
// @Success 200
// @Router /deployments/{id} [delete]
func (dh *DeploymentHandler) DeleteDeployment(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	deploymentId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = dh.deploymentApplication.DeleteDeployment(
		c.Request.Context(),
		deploymentId,
		currentUser.Id,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}

// DeployIngress FindAll godoc
// @Summary Deploy ingress
// @Tags deployment
// @ID deployIngress
// @Param X-User-ID header string true "User ID"
// @Param data body request.DeployIngress true "Deploy Ingress"
// @Product JSON
// @Success 200
// @Router /deployments/ingresses [post]
func (dh *DeploymentHandler) DeployIngress(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	var body request.DeployIngress
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := dh.deploymentApplication.DeployIngress(
		c.Request.Context(),
		currentUser.Id,
		body.EnvironmentId,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}
