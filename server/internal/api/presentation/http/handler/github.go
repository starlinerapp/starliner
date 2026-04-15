package handler

import (
	"net/http"
	"strconv"

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
// @Param organizationId path int true "Organization ID"
// @Product JSON
// @Success 200 {array} response.Repository
// @Router /github/repositories/{organizationId} [get]
func (gh *GithubHandler) GetRepositories(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("organizationId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repos, err := gh.githubApplication.GetRepositories(c.Request.Context(), currentUser.Id, organizationId)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewRepositories(repos))
}

// GetAllRepositories godoc
// @Summary Get All Repositories (owner only, unfiltered)
// @Tags github
// @ID getAllRepositories
// @Param X-User-ID header string true "User ID"
// @Param organizationId path int true "Organization ID"
// @Product JSON
// @Success 200 {array} response.Repository
// @Router /github/all-repositories/{organizationId} [get]
func (gh *GithubHandler) GetAllRepositories(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("organizationId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repos, err := gh.githubApplication.GetAllRepositories(c.Request.Context(), currentUser.Id, organizationId)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewRepositories(repos))
}

// GetRepositoryContents FindAll godoc
// @Summary Get Repository Content
// @Tags github
// @ID getRepositoryContents
// @Param X-User-ID header string true "User ID"
// @Param organizationId path int true "Organization ID"
// @Param owner path string true "Repository owner (user or org)"
// @Param repository path string true "Repository name"
// @Param path query string false "Path within the repository (e.g., src or src/main.go)"
// @Product JSON
// @Success 200 {array} response.RepositoryFile
// @Router /github/repositories/{organizationId}/{owner}/{repository}/contents [get]
func (gh *GithubHandler) GetRepositoryContents(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("organizationId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	owner := c.Param("owner")
	repository := c.Param("repository")
	repositoryPath := c.Query("path")

	repoContent, err := gh.githubApplication.GetRepositoryContents(
		c.Request.Context(),
		currentUser.Id,
		organizationId,
		owner,
		repository,
		repositoryPath,
	)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewRepositoryFiles(repoContent))
}
