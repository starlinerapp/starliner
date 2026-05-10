package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
	"starliner.app/internal/api/presentation/http/dto/response"
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusCreated, response.NewTeam(newTeam))
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignRepoToTeam godoc
// @Summary Assign a GitHub repository to a team
// @State core
// @Tags team
// @ID assignRepoToTeam
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Param data body request.AssignRepoToTeam true "Assign Repo"
// @Success 201
// @Router /teams/{teamId}/repos [post]
func (th *TeamHandler) AssignRepoToTeam(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var body request.AssignRepoToTeam
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = th.teamApplication.AssignRepoToTeam(c, currentUser.Id, teamId, body.GithubRepoId, body.RepoName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Status(http.StatusCreated)
}

// UnassignRepoFromTeam godoc
// @Summary Unassign a GitHub repository from a team
// @State core
// @Tags team
// @ID unassignRepoFromTeam
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Param repoId path int true "GitHub Repo ID"
// @Success 204
// @Router /teams/{teamId}/repos/{repoId} [delete]
func (th *TeamHandler) UnassignRepoFromTeam(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	repoId, err := strconv.ParseInt(c.Param("repoId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = th.teamApplication.UnassignRepoFromTeam(c, currentUser.Id, teamId, repoId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, response.NewTeamClusters(clusters))
}

// AssignClusterToTeam godoc
// @Summary Assign a cluster to a team
// @State core
// @Tags team
// @ID assignClusterToTeam
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Param clusterId path int true "Cluster ID"
// @Success 201
// @Router /teams/{teamId}/clusters/{clusterId} [post]
func (th *TeamHandler) AssignClusterToTeam(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	clusterId, err := strconv.ParseInt(c.Param("clusterId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = th.teamApplication.AssignClusterToTeam(c, currentUser.Id, teamId, clusterId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Status(http.StatusCreated)
}

// UnassignClusterFromTeam godoc
// @Summary Unassign a cluster from a team
// @State core
// @Tags team
// @ID unassignClusterFromTeam
// @Param X-User-ID header string true "User ID"
// @Param teamId path int true "Team ID"
// @Param clusterId path int true "Cluster ID"
// @Success 204
// @Router /teams/{teamId}/clusters/{clusterId} [delete]
func (th *TeamHandler) UnassignClusterFromTeam(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)

	teamId, err := strconv.ParseInt(c.Param("teamId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	clusterId, err := strconv.ParseInt(c.Param("clusterId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = th.teamApplication.UnassignClusterFromTeam(c, currentUser.Id, teamId, clusterId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Status(http.StatusNoContent)
}
