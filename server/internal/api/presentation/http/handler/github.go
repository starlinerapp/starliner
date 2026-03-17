package handler

import (
	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
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
// @Tags GitHub
// @ID getRepositories
// @Param X-User-ID header string true "User ID"
// @Product JSON
// @Success 200
// @Router /repositories [get]
func (gh *GithubHandler) GetRepositories(c *gin.Context) {
	_ = c.MustGet("user").(*value.User)

	gh.githubApplication.GetRepositories()
}
