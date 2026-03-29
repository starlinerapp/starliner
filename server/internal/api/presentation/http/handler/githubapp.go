package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
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
