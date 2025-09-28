package response

type Project struct {
	Id   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
