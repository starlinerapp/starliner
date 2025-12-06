package response

type Organization struct {
	Id      int64  `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Slug    string `json:"slug" binding:"required"`
	OwnerId int64  `json:"owner_id" binding:"required"`
}
