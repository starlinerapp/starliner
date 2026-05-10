package conf

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
}

func LoadConfig() (*Config, error) {
	var config Config

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfgType := reflect.TypeOf(config)
	for i := 0; i < cfgType.NumField(); i++ {
		if env := cfgType.Field(i).Tag.Get("mapstructure"); env != "" {
			if err := viper.BindEnv(env); err != nil {
				return nil, err
			}
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return &config, err
	}

	if err := validator.New().Struct(config); err != nil {
		return &config, err
	}

	return &config, nil
}
