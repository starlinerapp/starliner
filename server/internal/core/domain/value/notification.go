package value

type NotificationKind string

const (
	NotificationKindCluster    NotificationKind = "cluster"
	NotificationKindDeployment NotificationKind = "deployment"
)

type Notification struct {
	Kind       NotificationKind `json:"kind"`
	Status     string           `json:"status"`
	Message    string           `json:"message"`
	ResourceId *int64           `json:"resourceId,omitempty"`
}
