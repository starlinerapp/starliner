package objectstore

import "go.uber.org/fx"

var Module = fx.Module(
	"objectstore",
	fx.Provide(Connect),
	fx.Invoke(CreateBucket),
)
