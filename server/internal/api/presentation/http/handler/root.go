package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/presentation/http/dto/response"
)

type RootHandler struct{}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

// GetRoot FindAll godoc
// @Summary Get root
// @Tags root
// @ID getRoot
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} response.Root
// @Router / [get]
func (rh *RootHandler) GetRoot(c *gin.Context) {
	res := response.Root{Message: "Hello World!"}
	c.JSON(http.StatusOK, res)
}

// TriggerError godoc
// @Summary Trigger a test error
// @Tags root
// @ID triggerError
// @Product JSON
// @Success 500
// @Router /webhooks/debug/sentry [get]
//
// TriggerError deliberately returns a 500 with an attached error so we can
// confirm Sentry reporting works end-to-end on staging. Safe to remove once
// the integration is verified.
func (rh *RootHandler) TriggerError(c *gin.Context) {
	_ = c.Error(errors.New("test error: sentry integration check"))
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}
