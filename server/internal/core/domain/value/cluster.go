package value

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

type ClusterProvisionedSuccess struct {
	ClusterId        int64  `json:"clusterId"`
	ProvisioningId   string `json:"provisioningId"`
	IPv4Address      string `json:"ipv4Address"`
	PublicKey        string `json:"publicKey"`
	PrivateKey       string `json:"privateKey"`
	KubeconfigBase64 string `json:"kubeconfigBase64"`
	Logs             string `json:"logs"`
}

type ClusterProvisionedFailure struct {
	ClusterId int64  `json:"clusterId"`
	Reason    string `json:"reason"`
}

type ClusterDeletedSuccess struct {
	ClusterId int64 `json:"clusterId"`
}

type ClusterDeletedFailure struct {
	ClusterId int64  `json:"clusterId"`
	Reason    string `json:"reason"`
}
