package redis

import (
	"go.uber.org/fx"
	"starliner.app/internal/builder/domain/port"
)

var Module = fx.Module(
	"builder-redis",
	fx.Provide(
		NewStore,
		func(s *Store) port.RunnerStore { return s },
		func(s *Store) port.BuildDispatchStore { return s },
	),
)
