package config

import "go.uber.org/fx"

func ProvideConfig() (*Config, error) {
	return LoadConfig()
}

var Module = fx.Module(
	"config",
	fx.Provide(ProvideConfig),
)
