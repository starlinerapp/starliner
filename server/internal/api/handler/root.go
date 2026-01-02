package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/api/dto/response"
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
