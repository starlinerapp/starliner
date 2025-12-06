package response

type Root struct {
	Message string `json:"message" binding:"required"`
}
