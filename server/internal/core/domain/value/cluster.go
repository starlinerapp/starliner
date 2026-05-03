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
