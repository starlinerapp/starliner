package domain

type Cluster struct {
	Id             int64
	Name           string
	IPv4Address    *string
	PublicKey      *string
	PrivateKeyRef  *string
	OrganizationId int64
}
