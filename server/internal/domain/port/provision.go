package port

import "context"

type Provision interface {
	ProvisionServer(ctx context.Context, name string, publicKey []byte) (provisioningId string, ip string, err error)
	DeleteServer(ctx context.Context, provisioningId string) error
}
