package port

type InviteData struct {
	OrganizationName string
	InviteLink       string
}

type ResetData struct {
	PasswordResetLink string
}

type VerifyData struct {
	VerificationLink string
}

type Email interface {
	SendInvite(to string, inviteData InviteData) error
	SendVerificationEmail(to string, verifyData VerifyData) error
	SendResetPassword(to string, resetData ResetData) error
}
