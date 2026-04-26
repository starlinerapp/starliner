package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/presentation/http/dto/request"
)

type InternalHandler struct {
	userApplication *application.UserApplication
}

func NewInternalHandler(userApplication *application.UserApplication) *InternalHandler {
	return &InternalHandler{userApplication: userApplication}
}

// SendVerificationEmail godoc
// @Summary Send email verification
// @Tags internal
// @ID sendVerificationEmail
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param data body request.SendVerificationEmailRequest true "Verification"
// @Success 204
// @Router /internal/send-verification-email [post]
func (ih *InternalHandler) SendVerificationEmail(c *gin.Context) {
	var body request.SendVerificationEmailRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ih.userApplication.SendVerificationEmail(
		c.Request.Context(),
		body.To,
		body.VerificationLink,
	); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusNoContent)
}
