package request

type CreateOrganization struct {
	Name string `json:"name" binding:"required"`
}

type UpsertHetznerCredential struct {
	ApiKey string `json:"apiKey" binding:"required"`
}

type AcceptInvite struct {
	RecipientEmail string `json:"recipientEmail" binding:"required"`
	InviteId       string `json:"inviteId" binding:"required"`
}

type SendInvite struct {
	ToEmails        []string `json:"toEmails" binding:"required,min=1,dive,email"`
	InviteUrlPrefix string   `json:"inviteUrlPrefix" binding:"required,url"`
	TeamID          *int64   `json:"teamId"`
}

type RemoveOrganizationMember struct {
	UserID int64 `json:"userId" binding:"required"`
}
