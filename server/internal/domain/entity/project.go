package entity

type Project struct {
	Id             int64
	Name           string
	OrganisationId int64
	ClusterId      *int64
}
type ProjectWithEnvironments struct {
	Id             int64
	Name           string
	Environments   []*Environment
	OrganizationId int64
	ClusterId      *int64
}
