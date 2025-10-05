package request

type CreateProject struct {
	Name           string `json:"name" binding:"required"`
	OrganizationId int64  `json:"organization_id" binding:"required"`
}

type GetProject struct {
	Id int64 `json:"id" binding:"required"`
}
