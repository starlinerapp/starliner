package value

type ProvisionCluster struct {
	Id                     int64
	Name                   string
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
}

type ClusterDeleted struct {
	Id int64
}
