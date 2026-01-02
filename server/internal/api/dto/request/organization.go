package request

type CreateOrganization struct {
	Name string `json:"name" binding:"required"`
}
