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

type GithubAppHandler struct {
	githubAppApplication *application.GitHubAppApplication
}

func NewGithubAppHandler(githubAppApplication *application.GitHubAppApplication) *GithubAppHandler {
	return &GithubAppHandler{
		githubAppApplication: githubAppApplication,
	}
}

// CreateGithubApp FindAll godoc
// @Summary Create GitHub App
// @State core
// @Tags githubapp
// @ID createGithubApp
// @Param X-User-ID header string true "User ID"
// @Product JSON
// @Param data body request.CreateGithubApp true "Create GitHub App"
// @Success 201
// @Router /githubapps [post]
func (gh *GithubAppHandler) CreateGithubApp(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	var githubApp request.CreateGithubApp
	if err := c.BindJSON(&githubApp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := gh.githubAppApplication.CreateGitHubApp(c.Request.Context(), currentUser.Id, githubApp.OrganizationId, githubApp.InstallationId)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusCreated)
}

// GetGithubApp FindAll godoc
// @Summary Get GitHub App
// @State core
// @Tags githubapp
// @ID getGithubApp
// @Param X-User-ID header string true "User ID"
// @Param organizationId path int true "Organization ID"
// @Product JSON
// @Success 200 {object} response.GithubApp
// @Router /githubapps/{organizationId} [get]
func (gh *GithubAppHandler) GetGithubApp(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("organizationId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ghApp, err := gh.githubAppApplication.GetGithubApp(c.Request.Context(), currentUser.Id, organizationId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	if ghApp == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "GitHub App not found"})
		return
	}

	c.JSON(http.StatusOK, response.NewGithubApp(ghApp))
}
