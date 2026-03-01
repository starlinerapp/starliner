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

type ClusterHandler struct {
	clusterApplication      *application.ClusterApplication
	organizationApplication *application.OrganizationApplication
}

func NewClusterHandler(clusterApplication *application.ClusterApplication, organizationApplication *application.OrganizationApplication) *ClusterHandler {
	return &ClusterHandler{
		clusterApplication:      clusterApplication,
		organizationApplication: organizationApplication,
	}
}

// CreateCluster FindAll godoc
// @Summary Create Cluster
// @Tags cluster
// @ID createCluster
// @Param X-User-ID header string true "User ID"
// @Param data body request.CreateCluster true "Create Cluster"
// @Product JSON
// @Success 200 {object} response.Cluster
// @Router /clusters [post]
func (ch *ClusterHandler) CreateCluster(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	var cluster request.CreateCluster
	if err := c.BindJSON(&cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	newCluster, err := ch.clusterApplication.CreateCluster(c.Request.Context(), currentUser.Id, cluster.Name, cluster.OrganizationID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.NewCluster(newCluster))
}

// GetCluster FindAll godoc
// @Summary Get Cluster
// @Tags cluster
// @ID getCluster
// @Param X-User-ID header string true "User ID"
// @Product JSON
// @Param id path int true "Cluster ID"
// @Success 200 {object} response.Cluster
// @Router /clusters/{id} [get]
func (ch *ClusterHandler) GetCluster(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	clusterId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	cluster, err := ch.clusterApplication.GetUserCluster(c.Request.Context(), currentUser.Id, clusterId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	if cluster == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}
	c.JSON(http.StatusOK, response.NewCluster(cluster))
}

// GetClusterPrivateKey FindAll godoc
// @Summary Get Cluster Private Name
// @Tags cluster
// @ID getClusterPrivateKey
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Cluster ID"
// @Product application/octet-stream
// @Success 200 {file} file
// @Router /clusters/{id}/private-key [get]
func (ch *ClusterHandler) GetClusterPrivateKey(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	clusterId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	file, err := ch.clusterApplication.GetClusterPrivateKey(c.Request.Context(), clusterId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=private-key.pem")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.Itoa(len(file)))

	c.Data(http.StatusOK, "application/octet-stream", file)
}

// DeleteCluster FindAll godoc
// @Summary Delete Cluster
// @Tags cluster
// @ID deleteCluster
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Cluster ID"
// @Product JSON
// @Success 200
// @Router /clusters/{id} [delete]
func (ch *ClusterHandler) DeleteCluster(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	clusterId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = ch.clusterApplication.DeleteCluster(c.Request.Context(), currentUser.Id, clusterId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}
