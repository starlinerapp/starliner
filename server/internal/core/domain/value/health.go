package value

type Health string

const (
	Healthy   Health = "healthy"
	Unhealthy Health = "unhealthy"
)

type HealthStatus struct {
	DeploymentId int64
	Health       Health
	Status       string
}
