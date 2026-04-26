package request

type SendVerificationEmailRequest struct {
	To               string `json:"to" binding:"required,email"`
	VerificationLink string `json:"verificationUrl" binding:"required,url"`
}
