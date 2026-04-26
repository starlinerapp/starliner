package entity

import "time"

type GithubApp struct {
	ID             int64
	InstallationID int64
	OrganizationID int64
	CreatedAt      time.Time
}
