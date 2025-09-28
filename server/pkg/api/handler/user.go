package handler

import (
	"github.com/gin-gonic/gin"
	"starliner.app/pkg/api/dto/response"
	"starliner.app/pkg/domain"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUser FindAll godoc
// @Summary Get current user
// @Tags user
// @ID getUser
// @Product JSON
// @Success 200 {object} response.User
// @Router /me [get]
func (uh *UserHandler) GetUser(c *gin.Context) {
	user := c.MustGet("user").(*domain.User)
	res := response.User{UserId: user.Id, BetterAuthId: user.BetterAuthId}
	c.JSON(200, res)
}
