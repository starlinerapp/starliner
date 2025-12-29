package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/pkg/api/dto/request"
	"starliner.app/pkg/domain"
	"starliner.app/pkg/service"
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
// @Success 200
// @Router /clusters [post]
func (ch *ClusterHandler) CreateCluster(c *gin.Context) {
	currentUser := c.MustGet("user").(*domain.User)
	var cluster request.CreateCluster
	if err := c.BindJSON(&cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	orgs, err := ch.organizationService.GetUserOrganizations(c, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	found := false
	for _, org := range orgs {
		if org.Id == cluster.OrganizationID {
			found = true
		}
	}

	if !found {
		c.AbortWithStatusJSON(404, gin.H{"error": "Organization not found"})
		return
	}

	err = ch.clusterService.CreateCluster(cluster.Name, cluster.OrganizationID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusOK)
}
