package handler

import (
	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/response"
)

type GithubHandler struct {
	githubApplication *application.GitHubApplication
}

func NewGithubHandler(githubApplication *application.GitHubApplication) *GithubHandler {
	return &GithubHandler{
		githubApplication: githubApplication,
	}
}

// GetRepositories FindAll godoc
// @Summary Get Repositories
// @Tags github
// @ID getRepositories
// @Param X-User-ID header string true "User ID"
// @Product JSON
// @Success 200 {array} response.Repository
// @Router /github/repositories [get]
func (gh *GithubHandler) GetRepositories(c *gin.Context) {
	_ = c.MustGet("user").(*value.User)

	repos, err := gh.githubApplication.GetRepositories(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, response.NewRepositories(repos))
}
