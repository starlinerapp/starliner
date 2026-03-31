package response

import (
	"starliner.app/internal/api/domain/value"
	"time"
)

type GithubApp struct {
	InstallationID int64     `json:"installation_id" binding:"required"`
	OrganizationID int64     `json:"organization_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at" binding:"required"`
}

func NewGithubApp(app *value.GithubApp) GithubApp {
	return GithubApp{
		InstallationID: app.InstallationID,
		OrganizationID: app.OrganizationID,
		CreatedAt:      app.CreatedAt,
	}
}
