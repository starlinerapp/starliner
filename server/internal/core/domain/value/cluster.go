package value

import "errors"

type ServerType string

const (
	ServerTypeCX23  ServerType = "cx23"
	ServerTypeCPX22 ServerType = "cpx22"
)

type ProvisionCluster struct {
	Id                     int64
	Name                   string
	ServerType             ServerType
	OrganizationName       string
	ProvisioningCredential string
}

type DeleteCluster struct {
	Id                     int64
	ProvisioningId         string
	ProvisioningCredential string
}

type ReconcileCluster struct {
	Id                     int64
	ProvisioningId         string
	ProvisioningCredential string
}

type ReconcileClusterRequest struct {
	ClusterId      int64
	OrganizationId int64
	ProvisioningId string
}

type ClusterCreated struct {
	Id               int64
	ProvisioningId   string
	IPv4Address      string
	PublicKey        string
	PrivateKey       string
	KubeconfigBase64 string
	Logs             string
}

type ClusterDeleted struct {
	Id int64
}

type ClusterProvisioningFailed struct {
	Id   int64
	Logs string
}

var ErrClusterUnreachable = errors.New("cluster unreachable")

func IsClusterUnreachable(err error) bool {
	return errors.Is(err, ErrClusterUnreachable)
}
