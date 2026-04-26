package port

type InviteData struct {
	OrganizationName string
	InviteLink       string
}

type Email interface {
	SendInvite(to string, inviteData InviteData) error
}
