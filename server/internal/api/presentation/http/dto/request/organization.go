package request

type CreateOrganization struct {
	Name string `json:"name" binding:"required"`
}

type UpsertHetznerCredential struct {
	ApiKey string `json:"apiKey" binding:"required"`
}
