package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
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

// GetOrganizationRunners godoc
// @Summary Get all registered runners for an organization
// @State core
// @Tags runner
// @ID getOrganizationRunners
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Organization ID"
// @Success 200 {array} response.Runner
// @Router /organizations/{id}/runners [get]
func (rh *RunnerHandler) GetOrganizationRunners(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	organizationId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	runners, err := rh.runnerApplication.GetOrganizationRunners(c.Request.Context(), organizationId, currentUser.Id)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, response.NewRunners(runners))
}

// RegisterRunner godoc
// @Summary Register self-hosted runner
// @State runner
// @Tags runner
// @ID registerRunner
// @Produce json
// @Param data body request.RegisterRunner true "Register Runner"
// @Success 201
// @Router /api/runners/register [post]
func (rh *RunnerHandler) RegisterRunner(c *gin.Context) {
	var body request.RegisterRunner
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := rh.runnerApplication.RegisterRunner(
		c.Request.Context(),
		body.Token,
		body.Name,
		body.Labels,
		body.MaxConcurrentJobs,
	)
	if err != nil {
		if errors.Is(err, value.ErrInvalidRunnerRegistrationToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Status(http.StatusCreated)
}
