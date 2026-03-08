package entity

import "time"

type Project struct {
	Id             int64
	Name           string
	Environments   []*Environment
	OrganizationId int64
	ClusterId      *int64
	CreatedAt      time.Time
}
