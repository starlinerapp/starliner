package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
	"starliner.app/internal/api/presentation/http/mapper"
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
		mapper.MapEnvVarsFromRequest(body.Envs),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}

// UpdateImageDeployment FindAll godoc
// @Summary Update image deployment
// @Tags deployment
// @ID updateImageDeployment
// @Param X-User-ID header string true "User ID"
// @Param deploymentId path int true "Deployment ID"
// @Param data body request.UpdateImage true "Update Image"
// @Product JSON
// @Success 200
// @Router /deployments/images/{deploymentId} [put]
func (dh *DeploymentHandler) UpdateImageDeployment(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	deploymentId, err := strconv.ParseInt(c.Param("deploymentId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var body request.UpdateImage
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = dh.deploymentApplication.UpdateImageDeployment(
		c.Request.Context(),
		currentUser.Id,
		deploymentId,
		body.EnvironmentId,
		body.ImageName,
		body.Tag,
		body.Port,
		mapper.MapEnvVarsFromRequest(body.Envs),
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
		body.ServiceName,
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
		mapper.MapHostsFromRequest(body.IngressHosts),
		currentUser.Id,
		body.EnvironmentId,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}

// UpdateIngressDeployment FindAll godoc
// @Summary Update ingress deployment
// @Tags deployment
// @ID updateIngressDeployment
// @Param X-User-ID header string true "User ID"
// @Param deploymentId path int true "Deployment ID"
// @Param data body request.UpdateIngress true "Update Ingress"
// @Product JSON
// @Success 200
// @Router /deployments/ingresses/{deploymentId} [put]
func (dh *DeploymentHandler) UpdateIngressDeployment(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	deploymentId, err := strconv.ParseInt(c.Param("deploymentId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var body request.UpdateIngress
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = dh.deploymentApplication.UpdateIngressDeployment(
		c.Request.Context(),
		currentUser.Id,
		body.EnvironmentId,
		deploymentId,
		mapper.MapHostsFromRequest(body.IngressHosts),
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}

// DeployFromGitRepository FindAll godoc
// @Summary Deploy from Git Repository
// @Tags deployment
// @ID deployFromGitRepository
// @Param X-User-ID header string true "User ID"
// @Param data body request.DeployFromGit true "Deploy from Git"
// @Product JSON
// @Success 200
// @Router /deployments/git [post]
func (dh *DeploymentHandler) DeployFromGitRepository(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	var body request.DeployFromGit
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := dh.deploymentApplication.DeployFromGit(
		c.Request.Context(),
		currentUser.Id,
		body.EnvironmentId,
		body.ServiceName,
		body.GitUrl,
		body.ProjectRepositoryPath,
		body.DockerfilePath,
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
