package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
	"starliner.app/internal/api/presentation/http/dto/response"
	"starliner.app/internal/api/presentation/http/mapper"
)

type TeamHandler struct {
	teamApplication *application.TeamApplication
}

func NewTeamHandler(teamApplication *application.TeamApplication) *TeamHandler {
	return &TeamHandler{
		teamApplication: teamApplication,
	}
}

// CreateTeam FindAll godoc
// @Summary Create team
// @State core
// @Tags team
// @ID createTeam
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param data body request.CreateTeam true "Team slug (lowercase, alphanumeric, hyphens only)"
// @Param id path int true "Organization ID"
// @Success 201 {object} response.Team
// @Router /organizations/{id}/teams [post]
func (th *TeamHandler) CreateTeam(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	var team request.CreateTeam
	if err := c.BindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTeam, err := th.teamApplication.CreateTeam(c, team.Slug, organizationId, currentUser.Id)
	if err != nil {
		RespondInternalError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.NewTeam(newTeam))
}

// DeleteTeam FindAll godoc
// @Summary Delete Team
// @State core
// @Tags team
// @ID deleteTeam
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Success 200
// @Router /teams/{teamId} [delete]
func (th *TeamHandler) DeleteTeam(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = th.teamApplication.DeleteTeam(c, currentUser.Id, teamId)
	if err != nil {
		RespondInternalError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// GetUserTeams FindAll godoc
// @Summary Get User Teams
// @State core
// @Tags team
// @ID getUserTeams
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Organization ID"
// @Success 200 {array} response.Team
// @Router /organizations/{id}/teams [get]
func (th *TeamHandler) GetUserTeams(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	teams, err := th.teamApplication.GetUserTeams(c, organizationId, currentUser.Id)
	if err != nil {
		RespondInternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewTeams(teams))
}

// GetTeamMembers FindAll godoc
// @Summary Get Team Members
// @State core
// @Tags team
// @ID getTeamMembers
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Success 200 {array} response.User
// @Router /teams/{teamId}/members [get]
func (th *TeamHandler) GetTeamMembers(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	teamMembers, err := th.teamApplication.GetTeamMembers(c, currentUser.Id, teamId)
	if err != nil {
		RespondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewUsers(teamMembers))
}

// JoinTeam FindAll godoc
// @Summary Join a team by slug
// @State core
// @Tags team
// @ID joinTeam
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Organization ID"
// @Param data body request.JoinTeam true "Join Team"
// @Success 201
// @Router /organizations/{id}/teams/join [post]
func (th *TeamHandler) JoinTeam(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var body request.JoinTeam
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = th.teamApplication.JoinTeam(c, organizationId, currentUser.Id, body.Slug)
	if err != nil {
		RespondInternalError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

// AddTeamMember FindAll godoc
// @Summary Add organization member to team
// @State core
// @Tags team
// @ID addTeamMember
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Param data body request.AddTeamMember true "ID of member to add"
// @Success 201
// @Router /teams/{teamId}/members [post]
func (th *TeamHandler) AddTeamMember(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var body request.AddTeamMember
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = th.teamApplication.AddTeamMember(c, body.UserID, teamId, currentUser.Id)
	if err != nil {
		RespondInternalError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

// RemoveTeamMember FindAll godoc
// @Summary Remove organization member from team
// @State core
// @Tags team
// @ID removeTeamMember
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Param data body request.RemoveTeamMember true "ID of the organization member to remove from the team"
// @Success 204
// @Router /teams/{teamId}/members [delete]
func (th *TeamHandler) RemoveTeamMember(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var body request.RemoveTeamMember
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = th.teamApplication.RemoveTeamMember(c, body.UserID, teamId, currentUser.Id)
	if err != nil {
		RespondInternalError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
// SetTeamRepositories godoc
// @Summary Set repositories assigned to a team
// @State core
// @Tags team
// @ID setTeamRepositories
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Param data body request.SetTeamRepositories true "Team Repositories"
// @Success 204
// @Router /teams/{teamId}/repos [put]
func (th *TeamHandler) SetTeamRepositories(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var body request.SetTeamRepositories
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repos := mapper.MapTeamReposFromRequest(teamId, body.Repositories)

	err = th.teamApplication.SetTeamRepositories(c, currentUser.Id, teamId, repos)
	if err != nil {
		RespondInternalError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetTeamRepositories godoc
// @Summary Get repositories assigned to a team
// @State core
// @Tags team
// @ID getTeamRepositories
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Success 200 {array} response.TeamRepo
// @Router /teams/{teamId}/repos [get]
func (th *TeamHandler) GetTeamRepositories(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	repos, err := th.teamApplication.GetTeamRepositories(c, currentUser.Id, teamId)
	if err != nil {
		RespondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewTeamRepos(repos))
}

// GetTeamClusters godoc
// @Summary Get clusters assigned to a team
// @State core
// @Tags team
// @ID getTeamClusters
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Success 200 {array} response.TeamCluster
// @Router /teams/{teamId}/clusters [get]
func (th *TeamHandler) GetTeamClusters(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	clusters, err := th.teamApplication.GetTeamClusters(c, currentUser.Id, teamId)
	if err != nil {
		RespondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewTeamClusters(clusters))
}

// SetTeamClusters godoc
// @Summary Set clusters assigned to a team
// @State core
// @Tags team
// @ID setTeamClusters
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Param data body request.SetTeamClusters true "Team Clusters"
// @Success 204
// @Router /teams/{teamId}/clusters [put]
func (th *TeamHandler) SetTeamClusters(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var body request.SetTeamClusters
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clusters := mapper.MapTeamClustersFromRequest(teamId, body.Clusters)

	err = th.teamApplication.SetTeamClusters(c, currentUser.Id, teamId, clusters)
	if err != nil {
		RespondInternalError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
