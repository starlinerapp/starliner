package handler

import (
	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/response"
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
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} response.User
// @Router /me [get]
func (uh *UserHandler) GetUser(c *gin.Context) {
	user := c.MustGet("user").(*value.User)
	res := response.User{UserId: user.Id, BetterAuthId: user.BetterAuthId}
	c.JSON(200, res)
}
