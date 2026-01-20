package entity

type ClusterStatus string

const (
	ClusterPending ClusterStatus = "pending"
	ClusterRunning ClusterStatus = "running"
	ClusterDeleted ClusterStatus = "deleted"
)

type Cluster struct {
	Id             int64
	Name           string
	Status         ClusterStatus
	IPv4Address    *string
	PublicKey      *string
	PrivateKey     *string
	ProvisioningId *string
	Kubeconfig     *string
	OrganizationId int64
}
