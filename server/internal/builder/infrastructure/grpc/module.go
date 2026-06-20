package grpc

import "go.uber.org/fx"

var Module = fx.Module(
	"builder-grpc-clients",
	fx.Provide(
		NewAuthClient,
	),
)
