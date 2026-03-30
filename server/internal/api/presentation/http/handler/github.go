package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/response"
	"strconv"
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
// @Param organizationId path int true "Organization ID"
// @Product JSON
// @Success 200 {array} response.Repository
// @Router /github/repositories/{organizationId} [get]
func (gh *GithubHandler) GetRepositories(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("organizationId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	repos, err := gh.githubApplication.GetRepositories(c.Request.Context(), currentUser.Id, organizationId)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err})
		return
	}

	c.JSON(200, response.NewRepositories(repos))
}
