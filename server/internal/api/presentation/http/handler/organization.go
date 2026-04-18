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
// @Success 201 {object} response.Organization
// @Router /organizations [post]
func (oh *OrganizationHandler) CreateOrganization(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	var org request.CreateOrganization
	if err := c.BindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newOrg, err := oh.organizationApplication.CreateOrganization(c, org.Name, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusCreated, response.NewOrganization(newOrg))
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

// GetUserProjects FindAll godoc
// @Summary Get Organization Projects
// @Tags organization
// @ID getUserProjects
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Organization ID"
// @Success 200 {array} response.Project
// @Router /organizations/{id}/projects [get]
func (oh *OrganizationHandler) GetUserProjects(c *gin.Context) {
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

// UpsertHetznerCredential FindAll godoc
// @Summary Upsert Hetzner Provisioning Credential
// @Tags organization
// @ID upsertHetznerCredential
// @Product JSON
// @Param data body request.UpsertHetznerCredential true "Upsert Hetzner Credential"
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Organization ID"
// @Success 200
// @Router /organizations/{id}/settings/credential/hetzner [post]
func (oh *OrganizationHandler) UpsertHetznerCredential(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	var credential request.UpsertHetznerCredential
	if err := c.BindJSON(&credential); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = oh.organizationApplication.UpsertHetznerCredential(c.Request.Context(), currentUser.Id, organizationId, credential.ApiKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.Status(http.StatusOK)
}

// GetHetznerCredential FindAll godoc
// @Summary Get Hetzner Provisioning Credential
// @Tags organization
// @ID getHetznerCredential
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Organization ID"
// @Success 200 {object} response.GetOrganizationProvisioningCredentialResponse
// @Router /organizations/{id}/settings/credential/hetzner [get]
func (oh *OrganizationHandler) GetHetznerCredential(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	credential, err := oh.organizationApplication.GetHetznerCredential(c.Request.Context(), currentUser.Id, organizationId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if credential == nil {
		c.JSON(http.StatusOK, response.GetOrganizationProvisioningCredentialResponse{
			Credential: nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.GetOrganizationProvisioningCredentialResponse{
		Credential: &response.OrganizationProvisioningCredential{
			Provider: string(credential.Provider),
			Secret:   credential.Secret,
		},
	})
}

// CreateInvite godoc
// @Summary Create organization invite
// @Tags organization
// @ID createOrganizationInvite
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Organization ID"
// @Success 201 {object} response.OrganizationInvite
// @Router /organizations/{id}/invites [post]
func (oh *OrganizationHandler) CreateInvite(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	invite, err := oh.organizationApplication.CreateInvite(c.Request.Context(), currentUser.Id, organizationId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusCreated, response.NewOrganizationInvite(invite))
}

// GetInviteDetails godoc
// @Summary Get organization invite details
// @Tags organization
// @ID getOrganizationInviteDetails
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param inviteId path string true "Invite ID"
// @Success 200 {object} response.OrganizationInvite
// @Router /invites/{inviteId} [get]
func (oh *OrganizationHandler) GetInviteDetails(c *gin.Context) {
	inviteId := c.Param("inviteId")

	invite, err := oh.organizationApplication.GetInviteDetails(c.Request.Context(), inviteId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.NewOrganizationInvite(invite))
}

// AcceptInvite godoc
// @Summary Accept organization invite
// @Tags organization
// @ID acceptOrganizationInvite
// @Param X-User-ID header string true "User ID"
// @Param data body request.AcceptInvite true "Accept Invite"
// @Success 200
// @Router /invites/accept [post]
func (oh *OrganizationHandler) AcceptInvite(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	var body request.AcceptInvite
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := oh.organizationApplication.AcceptInvite(c.Request.Context(), body.InviteId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusOK)
}

// GetOrganizationMembers godoc
// @Summary Get all organization members
// @Tags organization
// @ID getOrganizationMembers
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Organization ID"
// @Success 200 {array} response.User
// @Router /organizations/{id}/members [get]
func (oh *OrganizationHandler) GetOrganizationMembers(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	members, err := oh.organizationApplication.GetOrganizationMembers(c.Request.Context(), currentUser.Id, organizationId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, response.NewUsers(members))
}
