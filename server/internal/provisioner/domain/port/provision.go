package port

import "context"

type Provision interface {
	ProvisionServer(ctx context.Context, provisioningCredential string, name string, publicKey []byte) (provisioningId string, ip string, err error)
	DeleteServer(ctx context.Context, provisioningCredential string, provisioningId string) error
}
