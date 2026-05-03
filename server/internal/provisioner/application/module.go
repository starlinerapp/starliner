package application

import (
	"go.uber.org/fx"
	"starliner.app/internal/provisioner/domain/port"
)

var Module = fx.Module(
	"application",
	fx.Provide(
		NewClusterLogApplication,
		func(ca *ClusterLogApplication) port.LogPublisher {
			return ca
		},
		NewClusterApplication,
		NewTTYApplication,
	),
)
