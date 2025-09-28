package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/pkg/api/dto/request"
	"starliner.app/pkg/api/dto/response"
	"starliner.app/pkg/domain"
	"starliner.app/pkg/service"
	"strconv"
)

type OrganizationHandler struct {
	organizationService *service.OrganizationService
}

func NewOrganizationHandler(organizationService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		organizationService: organizationService,
	}
}

// CreateOrganization FindAll godoc
// @Summary Create organization
// @Tags organization
// @ID createOrganization
// @Product JSON
// @Param data body request.CreateOrganization true "Create Organization"
// @Success 200
// @Router /organizations [post]
func (oh *OrganizationHandler) CreateOrganization(c *gin.Context) {
	currentUser := c.MustGet("user").(*domain.User)
	var org request.CreateOrganization
	if err := c.BindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_, err := oh.organizationService.CreateOrganization(c, org.Name, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusCreated)
}

// GetUserOrganizations FindAll godoc
// @Summary Get user organizations
// @Tags organization
// @ID getUserOrganizations
// @Product JSON
// @Success 200 {object} []response.Organization
// @Router /organizations [get]
func (oh *OrganizationHandler) GetUserOrganizations(c *gin.Context) {
	currentUser := c.MustGet("user").(*domain.User)
	organizations, err := oh.organizationService.GetUserOrganizations(c, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	res := make([]response.Organization, len(organizations))
	for i, org := range organizations {
		res[i] = response.Organization{
			Id:      org.Id,
			Name:    org.Name,
			Slug:    org.Slug,
			OwnerId: org.OwnerId,
		}
	}

	c.JSON(http.StatusOK, res)
}

// GetOrganizationProjects FindAll godoc
// @Summary Get Organization Projects
// @Tags organization
// @ID getOrganizationProjects
// @Product JSON
// @Param id path int true "Organization ID"
// @Success 200 {array} response.Project
// @Router /organizations/{id}/projects [get]
func (oh *OrganizationHandler) GetOrganizationProjects(c *gin.Context) {
	currentUser := c.MustGet("user").(*domain.User)
	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	projects, err := oh.organizationService.GetProjectsForUser(c, currentUser.Id, organizationId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	res := make([]response.Project, len(projects))
	for i, project := range projects {
		res[i] = response.Project{
			Id:   project.Id,
			Name: project.Name,
		}
	}
	c.JSON(http.StatusOK, res)
}
