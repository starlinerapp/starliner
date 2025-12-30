package domain

type Cluster struct {
	Id             int64
	Name           string
	IPv4Address    *string
	PublicKey      *string
	PrivateKey     *string
	OrganizationId int64
}
