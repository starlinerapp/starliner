package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RootHandler struct{}

type RootResponse struct {
	Message string `json:"message"`
}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

// GetRoot FindAll godoc
// @Summary Get Root
// @Tags root
// @ID GetRoot
// @Product JSON
// @Success 200 {object} RootResponse
// @Router / [get]
func (rh *RootHandler) GetRoot(c *gin.Context) {
	response := RootResponse{Message: "Hello World!"}
	c.JSON(http.StatusOK, response)
}
