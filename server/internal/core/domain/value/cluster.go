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
	CorrelationId          string
}

type DeleteCluster struct {
	Id                     int64
	ProvisioningId         string
	ProvisioningCredential string
	CorrelationId          string
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

type ClusterProvisionedSuccess struct {
	ClusterId        int64  `json:"clusterId"`
	ProvisioningId   string `json:"provisioningId"`
	IPv4Address      string `json:"ipv4Address"`
	PublicKey        string `json:"publicKey"`
	PrivateKey       string `json:"privateKey"`
	KubeconfigBase64 string `json:"kubeconfigBase64"`
	Logs             string `json:"logs"`
	CorrelationId    string `json:"correlationId"`
}

type ClusterProvisionedFailure struct {
	ClusterId     int64  `json:"clusterId"`
	Reason        string `json:"reason"`
	CorrelationId string `json:"correlationId"`
}

type ClusterDeletedSuccess struct {
	ClusterId     int64  `json:"clusterId"`
	CorrelationId string `json:"correlationId"`
}

type ClusterDeletedFailure struct {
	ClusterId     int64  `json:"clusterId"`
	Reason        string `json:"reason"`
	CorrelationId string `json:"correlationId"`
}

var ErrClusterUnreachable = errors.New("cluster unreachable")

func IsClusterUnreachable(err error) bool {
	return errors.Is(err, ErrClusterUnreachable)
}
