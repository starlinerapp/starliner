package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/response"
)

type RunnerHandler struct {
	runnerApplication *application.RunnerApplication
}

func NewRunnerHandler(runnerApplication *application.RunnerApplication) *RunnerHandler {
	return &RunnerHandler{
		runnerApplication: runnerApplication,
	}
}

// CreateRunner godoc
// @Summary Create self-hosted runner registration
// @State core
// @Tags runner
// @ID createRunner
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Organization ID"
// @Success 201 {object} response.CreateRunner
// @Router /organizations/{id}/runners [post]
func (rh *RunnerHandler) CreateRunner(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := rh.runnerApplication.CreateRunner(c, organizationId, currentUser.Id)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, response.NewCreateRunner(result))
}
