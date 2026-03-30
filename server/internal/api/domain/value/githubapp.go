package value

import "time"

type GithubApp struct {
	InstallationID int64
	OrganizationID int64
	CreatedAt      time.Time
}
