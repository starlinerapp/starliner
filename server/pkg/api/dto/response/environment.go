package response

type Environment struct {
	Id   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}
