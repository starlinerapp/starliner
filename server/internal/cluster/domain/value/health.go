package value

type Health string

const (
	Healthy   Health = "healthy"
	Unhealthy Health = "unhealthy"
)

type HealthStatus struct {
	Health Health
	Status string
}
