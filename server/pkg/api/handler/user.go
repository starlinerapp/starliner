package handler

import (
	"github.com/gin-gonic/gin"
	"starliner.app/pkg/domain"
)

type UserHandler struct{}

type UserResponse struct {
	UserId       int64  `json:"user_id"`
	BetterAuthId string `json:"better_auth_id"`
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUser FindAll godoc
// @Summary Get user
// @Tags user
// @ID getUser
// @Product JSON
// @Success 200 {object} UserResponse
// @Router /me [get]
func (uh *UserHandler) GetUser(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	response := UserResponse{UserId: user.Id, BetterAuthId: user.BetterAuthId}
	c.JSON(200, response)
}
