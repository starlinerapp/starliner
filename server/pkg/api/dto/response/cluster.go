package response

type Cluster struct {
	Id             int64   `json:"id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	IPv4Address    *string `json:"ipv4Address" binding:"required"`
	PublicKey      *string `json:"publicKey" binding:"required"`
	PrivateKey     *string `json:"privateKey" binding:"required"`
	OrganizationId int64   `json:"organizationId" binding:"required"`
}
