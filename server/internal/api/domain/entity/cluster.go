package entity

import "time"

type ClusterStatus string

const (
	ClusterPending ClusterStatus = "pending"
	ClusterRunning ClusterStatus = "running"
	ClusterDeleted ClusterStatus = "deleted"
)

type ServerType string

const (
	ServerTypeCX23  ServerType = "cx23"
	ServerTypeCCX33 ServerType = "ccx33"
)

type Cluster struct {
	Id             int64
	Name           string
	ServerType     ServerType
	Status         ClusterStatus
	User           string
	IPv4Address    *string
	PublicKey      *string
	PrivateKey     *string
	ProvisioningId *string
	Kubeconfig     *string
	OrganizationId int64
	CreatedAt      time.Time
}
