package entity

type Project struct {
	Id             int64
	Name           string
	Environments   []*Environment
	OrganizationId int64
	ClusterId      *int64
}
