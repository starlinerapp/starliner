package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/application"
	"starliner.app/internal/domain/value"
	"starliner.app/internal/presentation/http/dto/request"
	"starliner.app/internal/presentation/http/dto/response"
	"strconv"
)

type OrganizationHandler struct {
	organizationApplication *application.OrganizationApplication
}

func NewOrganizationHandler(organizationApplication *application.OrganizationApplication) *OrganizationHandler {
	return &OrganizationHandler{
		organizationApplication: organizationApplication,
	}
}

// CreateOrganization FindAll godoc
// @Summary Create organization
// @Tags organization
// @ID createOrganization
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param data body request.CreateOrganization true "Create Organization"
// @Success 201
// @Router /organizations [post]
func (oh *OrganizationHandler) CreateOrganization(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	var org request.CreateOrganization
	if err := c.BindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := oh.organizationApplication.CreateOrganization(c, org.Name, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusCreated)
}

// GetUserOrganizations FindAll godoc
// @Summary Get user organizations
// @Tags organization
// @ID getUserOrganizations
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} []response.Organization
// @Router /organizations [get]
func (oh *OrganizationHandler) GetUserOrganizations(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizations, err := oh.organizationApplication.GetUserOrganizations(c.Request.Context(), currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.NewOrganizations(organizations))
}

// GetOrganizationProjects FindAll godoc
// @Summary Get Organization Projects
// @Tags organization
// @ID getOrganizationProjects
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Organization ID"
// @Success 200 {array} response.Project
// @Router /organizations/{id}/projects [get]
func (oh *OrganizationHandler) GetOrganizationProjects(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	projects, err := oh.organizationApplication.GetProjectsForUser(c.Request.Context(), currentUser.Id, organizationId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, response.NewProjects(projects))
}

// GetOrganizationClusters FindAll godoc
// @Summary Get Organization Clusters
// @Tags organization
// @ID getOrganizationClusters
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Organization ID"
// @Success 200 {array} response.Cluster
// @Router /organizations/{id}/clusters [get]
func (oh *OrganizationHandler) GetOrganizationClusters(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	clusters, err := oh.organizationApplication.GetClustersForUser(c.Request.Context(), currentUser.Id, organizationId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, response.NewClusters(clusters))
}
