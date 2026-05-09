package conf

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
	fx.Provide(
		func() *Config {
			cfg, err := LoadConfig()
			if err != nil {
				panic(err)
			}
			return cfg
		},
	),
)
