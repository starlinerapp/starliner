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
	ToEmail         string `json:"toEmail" binding:"required,email"`
	InviteUrlPrefix string `json:"inviteUrlPrefix" binding:"required,url"`
}
