package model

type Project struct {
	Id             int64
	Name           string
	Environments   []Environment
	OrganizationId int64
}
