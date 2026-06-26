package entity

import "time"

type Runner struct {
	Id                int64
	OrganizationId    *int64
	Name              *string
	Status            string
	Labels            []string
	MaxConcurrentJobs int32
	DisabledAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
