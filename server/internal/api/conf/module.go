package conf

import (
	"go.uber.org/fx"
	"starliner.app/internal/core/conf"
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
		func(cfg *Config) conf.S3Config {
			return cfg
		},
		func(cfg *Config) conf.NatsConfig {
			return cfg
		},
		func(cfg *Config) conf.CryptoConfig {
			return cfg
		},
	),
)
