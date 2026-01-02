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

type ClusterHandler struct {
	clusterService      *service.ClusterService
	organizationService *service.OrganizationService
}

func NewClusterHandler(clusterService *service.ClusterService, organizationService *service.OrganizationService) *ClusterHandler {
	return &ClusterHandler{
		clusterService:      clusterService,
		organizationService: organizationService,
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
	currentUser := c.MustGet("user").(*model.User)
	var cluster request.CreateCluster
	if err := c.BindJSON(&cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	orgs, err := ch.organizationService.GetUserOrganizations(c.Request.Context(), currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	found := false
	for _, org := range orgs {
		if org.Id == cluster.OrganizationID {
			found = true
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	newCluster, err := ch.clusterService.CreateCluster(c, cluster.Name, cluster.OrganizationID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, &response.Cluster{
		Id:             newCluster.Id,
		Name:           newCluster.Name,
		Status:         newCluster.Status,
		IPv4Address:    newCluster.IPv4Address,
		OrganizationId: newCluster.OrganizationId,
	})
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
	currentUser := c.MustGet("user").(*model.User)
	clusterId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	cluster, err := ch.clusterService.GetUserCluster(c.Request.Context(), currentUser.Id, clusterId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	if cluster == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	c.JSON(http.StatusOK, &response.Cluster{
		Id:             cluster.Id,
		Name:           cluster.Name,
		Status:         cluster.Status,
		IPv4Address:    cluster.IPv4Address,
		OrganizationId: cluster.OrganizationId,
	})
}

// GetClusterPrivateKey FindAll godoc
// @Summary Get Cluster Private Key
// @Tags cluster
// @ID getClusterPrivateKey
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Cluster ID"
// @Product application/octet-stream
// @Success 200 {file} file
// @Router /clusters/{id}/private-key [get]
func (ch *ClusterHandler) GetClusterPrivateKey(c *gin.Context) {
	currentUser := c.MustGet("user").(*model.User)
	clusterId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	file, err := ch.clusterService.GetClusterPrivateKey(c.Request.Context(), clusterId, currentUser.Id)
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
	currentUser := c.MustGet("user").(*model.User)
	clusterId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = ch.clusterService.DeleteCluster(c.Request.Context(), currentUser.Id, clusterId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}
