package value

type ProvisioningCredentialProvider string

const (
	HetznerCredential ProvisioningCredentialProvider = "hetzner"
)

type ProvisioningCredential struct {
	Provider ProvisioningCredentialProvider
	Secret   string
}
