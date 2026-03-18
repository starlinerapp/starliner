package entity

import "time"

type Organization struct {
	Id      int64
	Name    string
	Slug    string
	OwnerId int64
}

type OrganizationInvite struct {
	Id               string
	OrganizationId   int64
	OrganizationName string
	ExpiresAt        time.Time
	CreatedAt        time.Time
}
