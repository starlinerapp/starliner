package port

import (
	"context"
	"starliner.app/internal/core/domain/value"
)

type Provision interface {
	ProvisionServer(ctx context.Context, clusterId int64, provisioningCredential string, name string, serverType value.ServerType, publicKey []byte) (provisioningId string, ip string, logs string, err error)
	DeleteServer(ctx context.Context, clusterId int64, provisioningCredential string, provisioningId string) error
	ReconcileServer(ctx context.Context, clusterId int64, provisioningCredential string, provisioningId string) (serverMissing bool, err error)
	DestroyServer(ctx context.Context, clusterId int64, provisioningCredential string, provisioningId string) error
}
