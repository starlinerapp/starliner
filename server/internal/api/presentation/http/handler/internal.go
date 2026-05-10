package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
// @State core
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

// SendResetPassword godoc
// @Summary Send password reset email
// @State core
// @Tags internal
// @ID sendResetPassword
// @Product JSON
// @Param X-User-ID header string true "User ID"
// @Param data body request.SendResetPasswordRequest true "Password reset"
// @Success 204
// @Router /internal/send-reset-password [post]
func (ih *InternalHandler) SendResetPassword(c *gin.Context) {
	var body request.SendResetPasswordRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ih.userApplication.SendResetPassword(
		c.Request.Context(),
		body.To,
		body.PasswordResetLink,
	); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.Status(http.StatusNoContent)
}
