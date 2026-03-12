package port

import (
	"context"
	"starliner.app/internal/core/domain/value"
)

type Provision interface {
	ProvisionServer(ctx context.Context, provisioningCredential string, name string, serverType value.ServerType, publicKey []byte) (provisioningId string, ip string, err error)
	DeleteServer(ctx context.Context, provisioningCredential string, provisioningId string) error
}
